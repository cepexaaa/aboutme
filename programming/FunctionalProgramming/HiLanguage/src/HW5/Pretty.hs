{-# LANGUAGE OverloadedStrings #-}

module HW5.Pretty
( prettyValue
, prettyError
, prettyExpr
, doc2String
) where

import qualified Data.ByteString               as BS
import           Data.Foldable                 (toList)
import qualified Data.Map                      as Map
import           Data.Ratio                    (denominator, numerator)
import           Data.Scientific               (FPFormat (Fixed),
                                                formatScientific,
                                                fromRationalRepetendUnlimited)
import           Data.Text                     ()
import qualified Data.Text                     as T
import           Data.Time                     (defaultTimeLocale, formatTime)
import           HW5.Base
import           Numeric                       (showHex)
import           Prettyprinter
import           Prettyprinter.Render.String   (renderString)
import           Prettyprinter.Render.Terminal

prettyFun :: HiFun -> String
prettyFun HiFunDiv            = "div"
prettyFun HiFunMul            = "mul"
prettyFun HiFunAdd            = "add"
prettyFun HiFunSub            = "sub"
prettyFun HiFunNot            = "not"
prettyFun HiFunAnd            = "and"
prettyFun HiFunOr             = "or"
prettyFun HiFunLessThan       = "less-than"
prettyFun HiFunGreaterThan    = "greater-than"
prettyFun HiFunEquals         = "equals"
prettyFun HiFunNotLessThan    = "not-less-than"
prettyFun HiFunNotGreaterThan = "not-greater-than"
prettyFun HiFunNotEquals      = "not-equals"
prettyFun HiFunIf             = "if"
prettyFun HiFunLength         = "length"
prettyFun HiFunToUpper        = "to-upper"
prettyFun HiFunToLower        = "to-lower"
prettyFun HiFunReverse        = "reverse"
prettyFun HiFunTrim           = "trim"
prettyFun HiFunList           = "list"
prettyFun HiFunRange          = "range"
prettyFun HiFunFold           = "fold"
prettyFun HiFunPackBytes      = "pack-bytes"
prettyFun HiFunUnpackBytes    = "unpack-bytes"
prettyFun HiFunEncodeUtf8     = "encode-utf8"
prettyFun HiFunDecodeUtf8     = "decode-utf8"
prettyFun HiFunZip            = "zip"
prettyFun HiFunUnzip          = "unzip"
prettyFun HiFunSerialise      = "serialise"
prettyFun HiFunDeserialise    = "deserialise"
prettyFun HiFunRead           = "read"
prettyFun HiFunWrite          = "write"
prettyFun HiFunMkDir          = "mkdir"
prettyFun HiFunChDir          = "cd"
prettyFun HiFunParseTime      = "parse-time"
prettyFun HiFunRand           = "rand"
prettyFun HiFunEcho           = "echo"
prettyFun HiFunCount          = "count"
prettyFun HiFunKeys           = "keys"
prettyFun HiFunValues         = "values"
prettyFun HiFunInvert         = "invert"


prettyRational :: Rational -> Doc AnsiStyle
prettyRational r
  | d == 1 = pretty n
  | repetend == Nothing = pretty (formatScientific Fixed Nothing sci)
  | otherwise =
      let whole = n `quot` d
          remainder = n `rem` d
      in if remainder == 0
           then pretty whole
           else if whole == 0
             then pretty n <> "/" <> pretty d
             else pretty whole <+> (if remainder >= 0 then "+" else "-")
                  <+> pretty (abs remainder) <> "/" <> pretty d
  where
    n = numerator r
    d = denominator r
    (sci, repetend) = fromRationalRepetendUnlimited r


prettyValue :: HiValue -> Doc AnsiStyle
prettyValue (HiValueNumber r) = prettyRational r
prettyValue (HiValueBool True) = "true"
prettyValue (HiValueBool False) = "false"
prettyValue HiValueNull = "null"
prettyValue (HiValueFunction f) = pretty (prettyFun f)
prettyValue (HiValueString s) = viaShow s
prettyValue (HiValueAction action) = prettyAction action
prettyValue (HiValueList xs) =
  let elements = map prettyValue (toList xs)
  in encloseSep "[ " " ]" ", " elements
prettyValue (HiValueBytes bs) =
  let hexBytes = BS.unpack bs
      formatByte b =
        let hex = showHex b ""
        in if length hex == 1 then "0" ++ hex else hex
      bytesStr = unwords (map formatByte hexBytes)
  in "[#" <+> pretty bytesStr <+> "#]"
prettyValue (HiValueTime time) =
  let timeStr = formatTime defaultTimeLocale "%Y-%m-%d %H:%M:%S%Q %Z" time
  in "parse-time(\"" <> pretty timeStr <> "\")"
prettyValue (HiValueDict dict) =
  let entries = Map.toList dict
      prettyEntry (k, v) = prettyValue k <> ": " <> prettyValue v
  in "{" <+> hsep (punctuate "," (map prettyEntry entries)) <+> "}"

doc2String :: Doc AnsiStyle -> String
doc2String doc = renderString (layoutPretty defaultLayoutOptions doc)


prettyExpr :: HiExpr -> Doc AnsiStyle
prettyExpr (HiExprValue v) = prettyValue v

prettyExpr (HiExprApply func args) =
  case func of
    HiExprValue (HiValueFunction f) -> pretty (prettyFun f)
    _                               -> prettyExpr func
  <> tupled1 (map prettyExpr args)
prettyExpr (HiExprRun expr) = prettyExpr expr <> "!"
prettyExpr (HiExprDict entries) =
  "{" <+> hsep (punctuate "," (map (\(k, v) -> prettyExpr k <> ":" <> prettyExpr v) entries)) <+> "}"


prettyError :: HiError -> Doc AnsiStyle
prettyError HiErrorInvalidArgument = "Invalid argument"
prettyError HiErrorInvalidFunction = "Invalid function"
prettyError HiErrorArityMismatch   = "Arity mismatch"
prettyError HiErrorDivideByZero    = "Division by zero"


tupled1 :: [Doc a] -> Doc a
tupled1 = encloseSep "(" ")" ", "

prettyAction :: HiAction -> Doc AnsiStyle
prettyAction (HiActionRead path) =
  "read" <> tupled1 [prettyValue (HiValueString (T.pack path))]
prettyAction (HiActionWrite path bs) =
  "write" <> tupled1 [prettyValue (HiValueString (T.pack path)), prettyValue (HiValueBytes bs)]
prettyAction (HiActionMkDir path) =
  "mkdir" <> tupled1 [prettyValue (HiValueString (T.pack path))]
prettyAction (HiActionChDir path) =
  "cd" <> tupled1 [prettyValue (HiValueString (T.pack path))]
prettyAction HiActionCwd = "cwd"
prettyAction HiActionNow = "now"
prettyAction (HiActionRand low high) = "rand" <> tupled1 [pretty low, pretty high]
prettyAction (HiActionEcho text) = "echo" <> tupled1 [prettyValue (HiValueString text)]
