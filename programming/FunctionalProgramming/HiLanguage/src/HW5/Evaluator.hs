{-# LANGUAGE FlexibleContexts #-}

module HW5.Evaluator (
 eval
) where

import qualified Codec.Compression.Zlib as Zlib
import           Codec.Serialise        (deserialise, serialise)
import           Control.Monad          (forM)
import           Control.Monad.Except
import           Control.Monad.Trans    (lift)
import qualified Data.ByteString        as BS
import qualified Data.ByteString.Lazy   as LBS
import           Data.Foldable          (foldl', toList)
import           Data.List              (sort)
import qualified Data.Map               as Map
import           Data.Ratio             (denominator, numerator, (%))
import           Data.Semigroup         (stimes)
import qualified Data.Sequence          as Seq
import qualified Data.Text              as T
import           Data.Text.Encoding     (decodeUtf8', encodeUtf8)
import           Data.Time              (NominalDiffTime, addUTCTime,
                                         diffUTCTime)
import           Data.Word              (Word8)
import           HW5.Base
import           Text.Read              (readMaybe)


eval :: HiMonad m => HiExpr -> m (Either HiError HiValue)
eval expr = runExceptT (evalExpr expr)

evalExpr :: HiMonad m => HiExpr -> ExceptT HiError m HiValue
evalExpr (HiExprRun expr) = do
  action <- evalExpr expr
  case action of
    HiValueAction act -> lift (runAction act)
    _                 -> throwError HiErrorInvalidArgument

evalExpr (HiExprValue v) = return v
evalExpr (HiExprApply funcExpr argsExprs) = do
  funcVal <- evalExpr funcExpr


  case (funcVal, argsExprs) of
    (HiValueFunction HiFunAnd, [a, b]) -> do
      aVal <- evalExpr a
      case aVal of
        HiValueBool False -> return aVal
        HiValueNull       -> return aVal
        _                 -> evalExpr b

    (HiValueFunction HiFunOr, [a, b]) -> do
      aVal <- evalExpr a
      case aVal of
        HiValueBool False -> evalExpr b
        HiValueNull       -> evalExpr b
        _                 -> return aVal

    (HiValueFunction HiFunIf, [cond, thenExpr, elseExpr]) -> do
      condVal <- evalExpr cond
      case condVal of
        HiValueBool True  -> evalExpr thenExpr
        HiValueBool False -> evalExpr elseExpr
        _                 -> throwError HiErrorInvalidArgument

    _ -> do
      argsVals <- mapM evalExpr argsExprs
      applyValue funcVal argsVals
evalExpr (HiExprDict entries) = do
  evaluatedEntries <- forM entries $ \(keyExpr, valueExpr) -> do
    key <- evalExpr keyExpr
    value <- evalExpr valueExpr
    pure (key, value)

  return (HiValueDict (Map.fromList evaluatedEntries))

applyFunction :: HiMonad m => HiFun -> [HiValue] -> ExceptT HiError m HiValue

applyFunction HiFunLength [HiValueString s] =
  return $ HiValueNumber (fromIntegral (T.length s))
applyFunction HiFunToUpper [HiValueString s] =
  return $ HiValueString (T.toUpper s)
applyFunction HiFunToLower [HiValueString s] =
  return $ HiValueString (T.toLower s)
applyFunction HiFunReverse [HiValueString s] =
  return $ HiValueString (T.reverse s)
applyFunction HiFunTrim [HiValueString s] =
  return $ HiValueString (T.strip s)


applyFunction HiFunEquals [HiValueString a, HiValueString b] =
  return $ HiValueBool (a == b)
applyFunction HiFunAdd [HiValueString a, HiValueString b] =
  return $ HiValueString (a <> b)
applyFunction HiFunMul [HiValueString s, HiValueNumber n]
  | denominator n == 1 && numerator n > 0 =
      return $ HiValueString (stimes (numerator n) s)
applyFunction HiFunMul [HiValueNumber n, HiValueString s]
  | denominator n == 1 && numerator n > 0 =
      return $ HiValueString (stimes (numerator n) s)
applyFunction HiFunDiv [HiValueString a, HiValueString b] =
  return $ HiValueString (a <> T.pack "/" <> b)


applyFunction HiFunList args =
  return $ HiValueList (Seq.fromList args)

applyFunction HiFunRange [HiValueNumber start, HiValueNumber end] = do
  let generate n
        | n <= end = HiValueNumber n : generate (n + 1)
        | otherwise = []
      nums = generate start
  return $ HiValueList (Seq.fromList nums)

applyFunction HiFunFold [HiValueFunction f, HiValueList xs] = do
  case Seq.viewl xs of
    Seq.EmptyL -> throwError HiErrorInvalidArgument
    first Seq.:< rest -> do
      let foldFunc acc x = do
            accVal <- acc
            applyFunction f [accVal, x]
      foldl' foldFunc (return first) rest


applyFunction HiFunLength [HiValueList xs] =
  return $ HiValueNumber (fromIntegral (Seq.length xs))

applyFunction HiFunReverse [HiValueList xs] =
  return $ HiValueList (Seq.reverse xs)

applyFunction HiFunAdd [HiValueList a, HiValueList b] =
  return $ HiValueList (a Seq.>< b)

applyFunction HiFunMul [HiValueList xs, HiValueNumber n]
  | denominator n == 1 && numerator n > 0 =
      return $ HiValueList (stimes (numerator n) xs)
applyFunction HiFunMul [HiValueNumber n, HiValueList xs]
  | denominator n == 1 && numerator n > 0 =
      return $ HiValueList (stimes (numerator n) xs)


applyFunction HiFunAdd [HiValueBytes a, HiValueBytes b] =
  return $ HiValueBytes (a <> b)

applyFunction HiFunMul [HiValueBytes bs, HiValueNumber n]
  | denominator n == 1 && numerator n > 0 =
      return $ HiValueBytes (stimes (numerator n) bs)
applyFunction HiFunMul [HiValueNumber n, HiValueBytes bs]
  | denominator n == 1 && numerator n > 0 =
      return $ HiValueBytes (stimes (numerator n) bs)

applyFunction HiFunParseTime [HiValueString timeStr] =
  case readMaybe (T.unpack timeStr) of
    Just t  -> return (HiValueTime t)
    Nothing -> return HiValueNull


applyFunction HiFunRand [HiValueNumber low, HiValueNumber high]
  | denominator low == 1 && denominator high == 1 = do
      let lowInt = fromInteger (numerator low) :: Int
          highInt = fromInteger (numerator high) :: Int
      if lowInt <= highInt
        then return (HiValueAction (HiActionRand lowInt highInt))
        else throwError HiErrorInvalidArgument

applyFunction HiFunRand _ =
  throwError HiErrorArityMismatch


applyFunction HiFunAdd [HiValueTime time, HiValueNumber seconds] = do
  let seconds' = fromRational seconds :: NominalDiffTime
  return (HiValueTime (addUTCTime seconds' time))

applyFunction HiFunAdd [HiValueNumber seconds, HiValueTime time] = do
  let seconds' = fromRational seconds :: NominalDiffTime
  return (HiValueTime (addUTCTime seconds' time))

applyFunction HiFunSub [HiValueTime time1, HiValueTime time2] = do
  let diff = diffUTCTime time1 time2 :: NominalDiffTime
  return (HiValueNumber (toRational diff))

applyFunction HiFunEcho [HiValueString text] =
  return (HiValueAction (HiActionEcho text))

applyFunction HiFunEcho _ =
  throwError HiErrorArityMismatch


applyFunction HiFunAdd args = arithmeticOp (+) args
applyFunction HiFunSub args = arithmeticOp (-) args
applyFunction HiFunMul args = arithmeticOp (*) args
applyFunction HiFunDiv args = arithmeticDivOp args


applyFunction HiFunNot [HiValueBool boolVal] = return $ HiValueBool (not boolVal)
applyFunction HiFunAnd [HiValueBool a, HiValueBool b] = return $ HiValueBool (a && b)
applyFunction HiFunOr [HiValueBool a, HiValueBool b] = return $ HiValueBool (a || b)

applyFunction HiFunLessThan [a, b] = return $ HiValueBool (a < b)
applyFunction HiFunGreaterThan [a, b] = return $ HiValueBool (a > b)
applyFunction HiFunEquals [a, b] = return $ HiValueBool (a == b)
applyFunction HiFunNotLessThan [a, b] = return $ HiValueBool (not (a < b))
applyFunction HiFunNotGreaterThan [a, b] = return $ HiValueBool (not (a > b))
applyFunction HiFunNotEquals [a, b] = return $ HiValueBool (a /= b)

applyFunction HiFunPackBytes [HiValueList xs] = do
  let toByte :: HiValue -> Maybe Word8
      toByte (HiValueNumber n) =
        if denominator n == 1 && numerator n >= 0 && numerator n <= 255
          then Just (fromIntegral (numerator n))
          else Nothing
      toByte _ = Nothing

  case mapM toByte (toList xs) of
    Just bytes -> return $ HiValueBytes (BS.pack bytes)
    Nothing    -> throwError HiErrorInvalidArgument

applyFunction HiFunUnpackBytes [HiValueBytes bs] = do
  let bytes = BS.unpack bs
      toNumber :: Word8 -> HiValue
      toNumber b = HiValueNumber (fromIntegral b)
  return $ HiValueList (Seq.fromList (map toNumber bytes))

applyFunction HiFunEncodeUtf8 [HiValueString s] =
  return $ HiValueBytes (encodeUtf8 s)

applyFunction HiFunDecodeUtf8 [HiValueBytes bs] =
  case decodeUtf8' bs of
    Right text -> return $ HiValueString text
    Left _     -> return HiValueNull

applyFunction HiFunZip [HiValueBytes bs] = do
  let lazyBs = LBS.fromStrict bs
      params = Zlib.defaultCompressParams {
        Zlib.compressLevel = Zlib.bestCompression
      }
      compressed = Zlib.compressWith params lazyBs
  return $ HiValueBytes (LBS.toStrict compressed)

applyFunction HiFunUnzip [HiValueBytes bs] = do
  let lazyBs = LBS.fromStrict bs
      decompressed = Zlib.decompress lazyBs
  return $ HiValueBytes (LBS.toStrict decompressed)

applyFunction HiFunSerialise [val] = do
  let lazyBytes = serialise val
  return $ HiValueBytes (LBS.toStrict lazyBytes)

applyFunction HiFunDeserialise [HiValueBytes bs] = do
  let lazyBs = LBS.fromStrict bs
  return $ deserialise lazyBs

applyFunction HiFunRead [HiValueString path] =
  return $ HiValueAction (HiActionRead (T.unpack path))

applyFunction HiFunWrite [HiValueString path, HiValueString content] =
  return $ HiValueAction (HiActionWrite (T.unpack path) (encodeUtf8 content))

applyFunction HiFunWrite [HiValueString path, HiValueBytes content] =
  return $ HiValueAction (HiActionWrite (T.unpack path) content)

applyFunction HiFunMkDir [HiValueString path] =
  return $ HiValueAction (HiActionMkDir (T.unpack path))

applyFunction HiFunChDir [HiValueString path] =
  return $ HiValueAction (HiActionChDir (T.unpack path))

applyFunction HiFunCount [arg] = do
  case arg of
    HiValueString s ->
      return $ countString s
    HiValueBytes bs ->
      return $ countBytes bs
    HiValueList xs ->
      return $ countList xs
    _ -> throwError HiErrorInvalidArgument
  where
    countString :: T.Text -> HiValue
    countString s =
      let freqMap = Map.fromListWith (+) [(T.singleton c, 1 :: Integer) | c <- T.unpack s]
      in HiValueDict (Map.fromList [(HiValueString k, HiValueNumber (v % 1))
                                   | (k, v) <- Map.toList freqMap])

    countBytes :: BS.ByteString -> HiValue
    countBytes bs =
      let freqMap = Map.fromListWith (+) [(toInteger b, 1 :: Integer) | b <- BS.unpack bs]
      in HiValueDict (Map.fromList [(HiValueNumber (k % 1),
                                    HiValueNumber (v % 1))
                                  | (k, v) <- Map.toList freqMap])

    countList :: Seq.Seq HiValue -> HiValue
    countList xs =
      let items = toList xs
          freqMap = foldl' incCount Map.empty items
      in HiValueDict freqMap
      where
        incCount m x = case Map.lookup x m of
          Nothing                -> Map.insert x (HiValueNumber (1 % 1)) m
          Just (HiValueNumber n) -> Map.insert x (HiValueNumber (n + 1 % 1)) m
          _                      -> m

applyFunction HiFunKeys [HiValueDict dict] =
  return $ HiValueList (Seq.fromList (Map.keys dict))

applyFunction HiFunValues [HiValueDict dict] =
  return $ HiValueList (Seq.fromList (Map.elems dict))

applyFunction HiFunInvert [HiValueDict dict] =
  return $ invertDict dict
  where
    invertDict :: Map.Map HiValue HiValue -> HiValue
    invertDict dict1 =
      let inverted = Map.foldlWithKey'
            (\m k v -> Map.insertWith (++) v [k] m)
            Map.empty
            dict1
          sorted = Map.map (HiValueList . Seq.fromList . sort) inverted
      in HiValueDict sorted

applyFunction HiFunIf args
  | length args == 3 = case args of
      [HiValueBool cond, thenVal, elseVal] ->
        return $ if cond then thenVal else elseVal
      _ -> throwError HiErrorInvalidArgument
  | otherwise = throwError HiErrorArityMismatch

applyFunction f args =
  case arityOf f of
    Variadic -> throwError HiErrorInvalidArgument
    Fixed n ->
      if length args == n
        then throwError HiErrorInvalidArgument
        else throwError HiErrorArityMismatch


data Arity = Fixed Int | Variadic

arityOf :: HiFun -> Arity
arityOf HiFunList           = Variadic
arityOf HiFunAdd            = Fixed 2
arityOf HiFunSub            = Fixed 2
arityOf HiFunMul            = Fixed 2
arityOf HiFunDiv            = Fixed 2
arityOf HiFunNot            = Fixed 1
arityOf HiFunAnd            = Fixed 2
arityOf HiFunOr             = Fixed 2
arityOf HiFunLessThan       = Fixed 2
arityOf HiFunGreaterThan    = Fixed 2
arityOf HiFunEquals         = Fixed 2
arityOf HiFunNotLessThan    = Fixed 2
arityOf HiFunNotGreaterThan = Fixed 2
arityOf HiFunNotEquals      = Fixed 2
arityOf HiFunIf             = Fixed 3
arityOf HiFunLength         = Fixed 1
arityOf HiFunToUpper        = Fixed 1
arityOf HiFunToLower        = Fixed 1
arityOf HiFunReverse        = Fixed 1
arityOf HiFunTrim           = Fixed 1
arityOf HiFunRange          = Fixed 2
arityOf HiFunFold           = Fixed 2
arityOf HiFunPackBytes      = Fixed 1
arityOf HiFunUnpackBytes    = Fixed 1
arityOf HiFunEncodeUtf8     = Fixed 1
arityOf HiFunDecodeUtf8     = Fixed 1
arityOf HiFunZip            = Fixed 1
arityOf HiFunUnzip          = Fixed 1
arityOf HiFunSerialise      = Fixed 1
arityOf HiFunDeserialise    = Fixed 1
arityOf HiFunRead           = Fixed 1
arityOf HiFunWrite          = Fixed 2
arityOf HiFunMkDir          = Fixed 1
arityOf HiFunChDir          = Fixed 1
arityOf HiFunParseTime      = Fixed 1
arityOf HiFunRand           = Fixed 2
arityOf HiFunEcho           = Fixed 1
arityOf HiFunCount          = Fixed 1
arityOf HiFunKeys           = Fixed 1
arityOf HiFunValues         = Fixed 1
arityOf HiFunInvert         = Fixed 1





applyValue :: HiMonad m => HiValue -> [HiValue] -> ExceptT HiError m HiValue
applyValue (HiValueFunction f) args = applyFunction f args
applyValue (HiValueString str) args = applyStringFunc str args
applyValue (HiValueBytes bs) args   = applyBytes bs args
applyValue (HiValueList xs) args    = applyList xs args
applyValue (HiValueDict dict) args  = applyDict dict args
applyValue _ _                      = throwError HiErrorInvalidFunction

applyDict :: HiMonad m => Map.Map HiValue HiValue -> [HiValue] -> ExceptT HiError m HiValue
applyDict dict [key] =
  case Map.lookup key dict of
    Just value -> return value
    Nothing    -> return HiValueNull
applyDict _ _ = throwError HiErrorArityMismatch


applyBytes :: HiMonad m => BS.ByteString -> [HiValue] -> ExceptT HiError m HiValue
applyBytes bs [HiValueNumber idx]
  | denominator idx == 1 =
      let i = fromInteger (numerator idx)
          len = BS.length bs
      in if i < 0 || i >= len
           then return HiValueNull
           else return $ HiValueNumber (fromIntegral (BS.index bs i))
  | otherwise = throwError HiErrorInvalidArgument
applyBytes bs [a, b] = do
  (l, r) <- bounds (BS.length bs) a b
  return $ HiValueBytes (BS.take (r - l) (BS.drop l bs))
  where
    toBound HiValueNull = return Nothing
    toBound (HiValueNumber x)
      | denominator x == 1 = return (Just (fromInteger (numerator x)))
      | otherwise = throwError HiErrorInvalidArgument
    toBound _ = throwError HiErrorInvalidArgument

    norm len i = if i < 0 then len + i else i
    clamp lo hi v = max lo (min hi v)

    bounds len a1 b1 = do
      s0 <- toBound a1
      e0 <- toBound b1
      let s1 = clamp 0 len (norm len (maybe 0 id s0))
          e1 = clamp 0 len (norm len (maybe len id e0))
          r  = if s1 <= e1 then e1 else s1
      return (s1, r)
applyBytes _ _ = throwError HiErrorArityMismatch



applyStringFunc :: HiMonad m => T.Text -> [HiValue] -> ExceptT HiError m HiValue
applyStringFunc s [HiValueNumber idx]
  | denominator idx == 1 =
      let i = fromInteger (numerator idx)
          len = T.length s
      in if i < 0 || i >= len
           then return HiValueNull
           else return $ HiValueString (T.singleton (T.index s i))
  | otherwise = throwError HiErrorInvalidArgument
applyStringFunc s [a, b] = do
  (l, r)<- bounds (T.length s) a b
  return $ HiValueString (T.take (r - l) (T.drop l s))
  where
    toBound HiValueNull = return Nothing
    toBound (HiValueNumber x)
      | denominator x == 1 = return (Just (fromInteger (numerator x)))
      | otherwise = throwError HiErrorInvalidArgument
    toBound _ = throwError HiErrorInvalidArgument

    norm len i = if i < 0 then len + i else i
    clamp lo hi v = max lo (min hi v)

    bounds len a1 b1 = do
      s0 <- toBound a1
      e0 <- toBound b1
      let s1 = clamp 0 len (norm len (maybe 0 id s0))
          e1 = clamp 0 len (norm len (maybe len id e0))
          r  = if s1 <= e1 then e1 else s1
      return (s1, r)
applyStringFunc _ _ = throwError HiErrorArityMismatch

applyList :: HiMonad m => Seq.Seq HiValue -> [HiValue] -> ExceptT HiError m HiValue
applyList xs [HiValueNumber idx]
  | denominator idx == 1 =
      let i = fromInteger (numerator idx)
          len = Seq.length xs
      in if i < 0 || i >= len
           then return HiValueNull
           else return (Seq.index xs i)
  | otherwise = throwError HiErrorInvalidArgument
applyList xs [a, b] = do
  (l, r) <- bounds (Seq.length xs) a b
  return $ HiValueList (Seq.take (r - l) (Seq.drop l xs))
  where
    toBound HiValueNull = return Nothing
    toBound (HiValueNumber x)
      | denominator x == 1 = return (Just (fromInteger (numerator x)))
      | otherwise = throwError HiErrorInvalidArgument
    toBound _ = throwError HiErrorInvalidArgument

    norm len i = if i < 0 then len + i else i
    clamp lo hi v = max lo (min hi v)

    bounds len a1 b1 = do
      s0 <- toBound a1
      e0 <- toBound b1
      let s1 = clamp 0 len (norm len (maybe 0 id s0))
          e1 = clamp 0 len (norm len (maybe len id e0))
          r  = if s1 <= e1 then e1 else s1
      return (s1, r)
applyList _ _ = throwError HiErrorArityMismatch


arithmeticOp :: HiMonad m => (Rational -> Rational -> Rational) -> [HiValue] -> ExceptT HiError m HiValue
arithmeticOp op [HiValueNumber a, HiValueNumber b] =
  return $ HiValueNumber (a `op` b)
arithmeticOp _ args = checkArithmeticArgs args

arithmeticDivOp :: HiMonad m => [HiValue] -> ExceptT HiError m HiValue
arithmeticDivOp [HiValueNumber a, HiValueNumber b]
  | b == 0 = throwError HiErrorDivideByZero
  | otherwise = return $ HiValueNumber (a / b)
arithmeticDivOp args = checkArithmeticArgs args

checkArithmeticArgs :: HiMonad m => [HiValue] -> ExceptT HiError m HiValue
checkArithmeticArgs args
  | length args /= 2 = throwError HiErrorArityMismatch
  | not (all isNumber args) = throwError HiErrorInvalidArgument
  | otherwise = error "Should not happen"
  where
    isNumber (HiValueNumber _) = True
    isNumber _                 = False

