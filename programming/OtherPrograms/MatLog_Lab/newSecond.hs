-- myparser.hs

{-# LANGUAGE LambdaCase #-}

module Main (main) where

import Data.Map
import Data.List
import Data.Char
import Control.Applicative
import System.IO( isEOF )

{-

Grammatic

S ::= E

E ::= D E'

E' ::= -> D E'
      | <eps>

D ::= C D'

D' ::= | C D'
      | <eps>

C ::= I C'

C' ::= & I C'
      | <eps>

I ::= ! I
      | ( E )
      | <Var>

-}

-- Expression --

data Expr
  = E_Var String
  | Expr :| Expr
  | Expr :& Expr
  | Expr :-> Expr
  | E_Inverse Expr
  deriving (Eq, Ord)

instance Show Expr where
    show = \case
      E_Var s -> id s
      lhs :|  rhs -> "(" ++ show lhs ++ "|" ++ show rhs ++ ")"
      lhs :&  rhs -> "(" ++ show lhs ++ "&" ++ show rhs ++ ")"
      lhs :-> rhs -> "(" ++ show lhs ++ "->" ++ show rhs ++ ")"
      E_Inverse e -> "(!" ++ show e ++ ")"

-- Lexing --

data Token 
    = T_Var String
    | T_Disj
    | T_Conj
    | T_Impl
    | T_Inverse
    | T_LPAR
    | T_RPAR
    | T_HALT_
    | T_Comma
    | T_Turn
    deriving (Eq, Show)

isUCAlpha :: Char -> Bool
isUCAlpha c = ord 'A' <= ord c && ord c <= ord 'Z'

isVarSymbol :: Char -> Bool
isVarSymbol c = isUCAlpha c || isDigit c || c == '\0039'

tokenize :: String -> [Token]
tokenize [] = []
tokenize ('|' : '-' : rest) = T_Turn    : tokenize rest
tokenize ('|' : rest)       = T_Disj    : tokenize rest
tokenize ('&' : rest)       = T_Conj    : tokenize rest
tokenize ('-' : '>' : rest) = T_Impl    : tokenize rest
tokenize ('!' : rest)       = T_Inverse : tokenize rest
tokenize ('(' : rest)       = T_LPAR    : tokenize rest
tokenize (')' : rest)       = T_RPAR    : tokenize rest
tokenize (',' : rest)       = T_Comma   : tokenize rest
tokenize s@(h : _)
    | isSpace h   = let(spaces, rest)     = span isSpace     s
                    in tokenize rest     
    | isUCAlpha h = let(varSymbols, rest) = span isVarSymbol s
                    in T_Var varSymbols : tokenize rest
tokenize _ = [T_HALT_]

-- Parser --

newtype Parser s a = Parser { runParser :: [s] -> Maybe (a, [s]) }

instance Functor (Parser s) where
    fmap f (Parser p) = Parser $ \ss -> case p ss of
        Nothing -> Nothing
        Just (x, ss') -> Just (f x, ss')

instance Applicative (Parser s) where
    pure x = Parser $ \ss -> Just (x, ss)

    Parser pf <*> Parser px = Parser $ \ss -> case pf ss of
        Nothing -> Nothing
        Just (f, ss') -> case px ss' of
            Nothing -> Nothing
            Just (x, ss'') -> Just (f x, ss'')

instance Monad (Parser s) where
    Parser p >>= f = Parser $ \ss -> case p ss of
        Nothing -> Nothing
        Just (x, ss') -> runParser (f x) ss'

instance Alternative (Parser s) where
    empty = Parser $ const Nothing

    Parser p1 <|> Parser p2 = Parser $ \ss -> case p1 ss of
        Nothing -> p2 ss
        def -> def

instance MonadFail (Parser s) where
    fail _ = Control.Applicative.empty

-- Parser tools --

ok :: Parser s ()
ok = pure ()

eof :: Parser s ()
eof = Parser $ \case
  [] -> Just ((), [])
  _ -> Nothing

pHead :: Parser s s
pHead = Parser $ \case
  [] -> Nothing
  (s : ss') -> Just (s, ss')

satisfies :: (s -> Bool) -> Parser s s
satisfies predicate = Parser $ \case
  [] -> Nothing
  (s : ss') -> if predicate s then Just (s, ss') else Nothing

element :: Eq s => s -> Parser s s
element s = satisfies (== s)

stream :: Eq s => [s] -> Parser s [s]
stream [] = pure []
stream (s : ss') = element s >>= \ps -> stream ss' >>= \pss' -> pure (ps : pss')

-- Expression parsing --

parseE :: Parser Token Expr
parseE     = (do d <- parseD
                 element T_Impl
                 e <- parseE
                 pure (d :-> e)) <|>
              parseD

parseD :: Parser Token Expr
parseD = parseC >>= parseD'

parseD' :: Expr -> Parser Token Expr
parseD' acc = (do element T_Disj
                  c <- parseC
                  parseD' (acc :| c)) <|>
              pure acc 

parseC :: Parser Token Expr
parseC = parseI >>= parseC'

parseC' :: Expr -> Parser Token Expr
parseC' acc = (do element T_Conj
                  i <- parseI
                  parseC' (acc :& i)) <|>
              pure acc 

parseI :: Parser Token Expr
parseI = (do element T_Inverse
             E_Inverse <$> parseI) <|>
         (do element T_LPAR
             e <- parseE
             element T_RPAR
             pure e) <|>
         (do T_Var s <- satisfies (\case { T_Var _ -> True; _ -> False })
             pure (E_Var s))

parseExpr :: String -> Maybe Expr
parseExpr s = fst <$> runParser (parseE <* eof) (tokenize s)

-- Metaexpression parsing --

-- metaexpression
type HypSet = Map Expr Int

type ContextSet = Map Expr Int

addElemToHypSet :: Expr -> HypSet -> HypSet
addElemToHypSet e hs = if member e hs 
                   then adjust ((+) 1) e hs
                   else Data.Map.insert e 1 hs

addListToHypSet :: [Expr] -> HypSet -> HypSet
addListToHypSet [] m = m
addListToHypSet (x:xs) m = addListToHypSet xs (addElemToHypSet x m)

addLeftImpls :: (ContextSet, Expr) -> (HypSet, MainExpr)
addLeftImpls (cs, (l :-> r)) = addLeftImpls ((addElemToHypSet l cs), r)
addLeftImpls hs = hs

type MainExpr = Expr

data Line = Line (ContextSet, HypSet, MainExpr) [Expr] Expr deriving Eq

showListOf :: [Expr] -> String
showListOf [] = ""
showListOf [h] = show h
showListOf (h:t) = show h ++ "," ++ showListOf t 

instance Show Line where
    show (Line cs c e) = showListOf c ++ "|-" ++ show e

getContext :: Line -> [Expr]
getContext (Line cs le e) = le

getContextSet :: Line -> ContextSet
getContextSet (Line (cs, hs, me) le e) = cs

getHypSet :: Line -> HypSet
getHypSet (Line (cs, hs, me) le e) = hs

getExpr :: Line -> Expr
getExpr (Line cs le e) = e

getME :: Line -> Expr
getME (Line (cs, hs, me) le e) = me

lineSets :: [Expr] -> Expr -> (ContextSet, HypSet, MainExpr)
lineSets le e = (\cs -> ((\sec -> (cs, fst sec, snd sec)) (addLeftImpls (cs, e)))) (addListToHypSet le Data.Map.empty)

-- -- -- -- -- --

parseContext :: Parser Token [Expr]
parseContext = (do e <- parseE 
                   t <- parseContext' [e]
                   pure t) <|>
               pure []

parseContext' :: [Expr] -> Parser Token [Expr]
parseContext' acc = (do element T_Comma
                        e <- parseE
                        parseContext' (acc ++ [e])) <|>
                    pure acc 

parseL :: Parser Token Line
parseL = (do c <- parseContext
             element T_Turn
             e <- parseE
             pure (Line (lineSets c e) c e))

parseLine :: String -> Maybe Line
parseLine s = fst <$> runParser (parseL <* eof) (tokenize s)

-- Tools

unJust = (\(Just e) -> e)

pack :: a -> [a]
pack = (\l -> [l])

(>>>) :: Eq a => Maybe a -> Maybe a -> Maybe a
(>>>) f s = if f == Nothing then s else f

atII :: [[a]] -> Int -> Int -> a
atII lla i1 i2 = lla !! i1 !! i2

toBool :: Maybe Bool -> Bool
toBool (Just True) = True
toBool _ = False

-- Deconstructors for expressions

headByImpl :: Expr -> Maybe Expr
headByImpl (a :-> b) = Just a
headByImpl _ = Nothing

tailByImpl :: Expr -> Maybe Expr
tailByImpl (a :-> b) = Just b
tailByImpl _ = Nothing

headByConj :: Expr -> Maybe Expr
headByConj (a :& b) = Just a
headByConj _ = Nothing

tailByConj :: Expr -> Maybe Expr
tailByConj (a :& b) = Just b
tailByConj _ = Nothing

headByDisj :: Expr -> Maybe Expr
headByDisj (a :| b) = Just a
headByDisj _ = Nothing

tailByDisj :: Expr -> Maybe Expr
tailByDisj (a :| b) = Just b
tailByDisj _ = Nothing

unInverse :: Expr -> Maybe Expr
unInverse (E_Inverse e) = Just e
unInverse _ = Nothing

splitToListByImpl :: Int -> Expr -> Maybe [Expr]
splitToListByImpl sch a
                   | sch == 1 = Just [a]
                   | sch > 1  = if headByImpl a == Nothing then Nothing
                                else ((++) [unJust $ headByImpl a] <$> (splitToListByImpl (sch - 1) $ unJust $ tailByImpl a))

splitByImpl :: [Int] -> Expr -> Maybe [[Expr]]
splitByImpl sch a
             | length sch == 1 = pack <$> splitToListByImpl (head sch) a
             | length sch > 1  = if headByImpl a == Nothing then Nothing
                                 else if (splitToListByImpl (head sch) (unJust $ headByImpl a)) == Nothing || (splitByImpl (tail sch) $ unJust $ tailByImpl a) == Nothing
                                      then Nothing
                                      else ((++) $ pack $ unJust $ (splitToListByImpl (head sch) (unJust $ headByImpl a))) <$> (splitByImpl (tail sch) $ unJust $ tailByImpl a)

-- Formalazing --

axSch :: Expr -> Maybe Int 
axSch me = ((splitByImpl [1,1,1] me) >>= (\lle -> if atII lle 0 0 == atII lle 2 0 then Just 1 else Nothing))

        >>>((splitByImpl [2,3,2] me) >>= (\lle -> if atII lle 0 0 == atII lle 1 0 &&
                                                     atII lle 1 0 == atII lle 2 0 &&
                                                     atII lle 0 1 == atII lle 1 1 && 
                                                     atII lle 1 2 == atII lle 2 1 then Just 2 else Nothing))

        >>>((splitByImpl [1,1,1] me) >>= (\lle -> if ((atII lle 0 0) :& (atII lle 1 0)) == atII lle 2 0 then Just 3 else Nothing))

        >>>((splitByImpl [1,1] me)   >>= (\lle -> if headByConj (atII lle 0 0) == Nothing 
                                                  then Nothing 
                                                  else if (unJust $ headByConj $ atII lle 0 0) == (atII lle 1 0) 
                                                       then Just 4 
                                                       else Nothing))

        >>>((splitByImpl [1,1] me)   >>= (\lle -> if headByConj (atII lle 0 0) == Nothing 
                                                  then Nothing
                                                  else if (unJust $ tailByConj $ atII lle 0 0) == (atII lle 1 0) 
                                                       then Just 5
                                                       else Nothing))

        >>>((splitByImpl [1,1] me)   >>= (\lle -> if headByDisj (atII lle 1 0) == Nothing 
                                                  then Nothing 
                                                  else if (atII lle 0 0) == (unJust $ headByDisj $ atII lle 1 0) 
                                                       then Just 6 
                                                       else Nothing))

        >>>((splitByImpl [1,1] me)   >>= (\lle -> if headByDisj (atII lle 1 0) == Nothing 
                                                  then Nothing 
                                                  else if (atII lle 0 0) == (unJust $ tailByDisj $ atII lle 1 0) 
                                                       then Just 7 
                                                       else Nothing))

        >>>((splitByImpl [2,2,2] me) >>= (\lle -> if atII lle 0 1 == atII lle 1 1 &&
                                                     atII lle 1 1 == atII lle 2 1 &&
                                                     ((atII lle 0 0) :| (atII lle 1 0)) == atII lle 2 0 then Just 8 else Nothing))

        >>>((splitByImpl [2,2,1] me) >>= (\lle -> if atII lle 0 0 == atII lle 1 0 &&
                                                     toBool ((==) (atII lle 0 0) <$> unInverse (atII lle 2 0)) &&
                                                     toBool ((==) (atII lle 0 1) <$> unInverse (atII lle 1 1)) then Just 9 else Nothing))

        >>>((splitByImpl [1,1] me)   >>= (\lle -> if toBool ((==) (atII lle 1 0) <$> (unInverse (atII lle 0 0) >>= unInverse)) then Just 10 else Nothing))

deducted :: [(Int, Line)] -> Line -> Maybe Int
deducted [] l = Nothing
deducted ll l = if (\pl -> (getHypSet l == getHypSet pl) && (getME l == getME pl)) (snd $ head ll)
                then Just $ fst $ head ll
                else deducted (tail ll) l

findLeftParts :: [(Int, Line)] -> Map Expr Int -> Line -> Map Expr Int
findLeftParts [] leftParts _ = leftParts
findLeftParts ll leftParts b = if (getContextSet b) == (getContextSet $ snd $ head ll)
                               then case getExpr $ snd $ head ll of
                                 (a :-> bb) -> if (getExpr b == bb) 
                                               then findLeftParts (tail ll) (Data.Map.insert a (fst $ head ll) leftParts) b
                                               else findLeftParts (tail ll) leftParts b
                                 _ -> findLeftParts (tail ll) leftParts b
                               else findLeftParts (tail ll) leftParts b

isAinMP :: [(Int, Line)] -> Map Expr Int -> Line -> Maybe (Int, Int)
isAinMP [] _ _ = Nothing
isAinMP ll leftParts b = if (getContextSet b == (getContextSet $ snd $ head ll)) && (member (getExpr $ snd $ head ll) leftParts)
                         then Just (fst $ head ll, unJust $ Data.Map.lookup (getExpr $ snd $ head ll) leftParts)
                         else isAinMP (tail ll) leftParts b

modusPonens :: [(Int, Line)] -> Line -> Maybe (Int, Int)
modusPonens ll l = isAinMP ll (findLeftParts ll Data.Map.empty l) l

fromHyp :: Line -> Maybe Int
fromHyp (Line _ le e) = ((+) 1) <$> elemIndex e le

comment :: [(Int, Line)] -> String
comment [] = undefined
comment ((num, l):ls) = case axSch $ getExpr l of
                          Just a -> "[Ax. sch. " ++ show a ++ "]"
                          Nothing -> case fromHyp l of
                                       Just b -> "[Hyp. " ++ show b ++ "]"
                                       Nothing -> case deducted ls l of
                                                    Just c -> let dedLine = snd $ ls !! (c - 1)
                                                              in if isIncorrect dedLine
                                                                 then "[Ded. " ++ show c ++ "; from Incorrect]"
                                                                 else "[Ded. " ++ show c ++ "]"
                                                    Nothing -> case modusPonens ls l of
                                                                 Just (a1, a2) -> let mpLine1 = snd $ ls !! (a1 - 1)
                                                                                      mpLine2 = snd $ ls !! (a2 - 1)
                                                                                  in if isIncorrect mpLine1 || isIncorrect mpLine2
                                                                                     then "[M.P. " ++ show a1 ++ ", " ++ show a2 ++ "; from Incorrect]"
                                                                                     else "[M.P. " ++ show a1 ++ ", " ++ show a2 ++ "]"
                                                                 Nothing -> "[Incorrect]"
  where
    isIncorrect line = case comment ((num, line):ls) of
                         "[Incorrect]" -> True
                         _ -> False

formalize :: [(Int, Line)] -> String
formalize ll = "[" ++ (show $ fst $ head ll) ++ "] " ++ (show $ snd $ head ll) ++ " " ++ (comment ll)

--добавить в список и дать представление переданной строки
processNewLine :: [(Int, Line)] -> String -> ([(Int, Line)], Maybe String)
processNewLine ll s = (\l -> if l == Nothing 
                             then (ll, Nothing) 
                             else (\new_ll -> (new_ll, Just $ formalize new_ll)) $ [(if ll == []
                                                                                     then 1
                                                                                     else (fst $ head ll) + 1, unJust l)] ++ ll) $ parseLine s

--прочитать строку, передать в обработку, если норм строка, то дальше 
ioLine :: [(Int, Line)] -> IO ()
ioLine ll = do
    end <- isEOF
    if end 
    then return ()
    else do s <- getLine
            (\p -> if snd p == Nothing 
                   then return ()
                   else do putStrLn (unJust $ snd p) 
                           ioLine $ fst p) (processNewLine ll s)

main :: IO ()
main = ioLine [] 