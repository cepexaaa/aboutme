{-# LANGUAGE OverloadedStrings #-}

module HW5.Parser
( parse
, exprParser
) where


import qualified Control.Monad.Combinators.Expr as Expr
import           Data.Bits                      (shiftL, (.|.))
import qualified Data.ByteString                as BS
import           Data.Char                      (isHexDigit, ord)
import           Data.Ratio                     ()
import           Data.Scientific                ()
import qualified Data.Text                      as T
import           Data.Void                      (Void)
import           Data.Word                      (Word8)
import           HW5.Base
import           Text.Megaparsec                hiding (parse)
import qualified Text.Megaparsec.Char           as C
import           Text.Megaparsec.Char
import qualified Text.Megaparsec.Char.Lexer     as L


type Parser = Parsec Void String


parse :: String -> Either (ParseErrorBundle String Void) HiExpr
parse = runParser (sc *> exprParser <* eof) ""


sc :: Parser ()
sc = L.space space1 empty empty

lexeme :: Parser a -> Parser a
lexeme = L.lexeme sc

symbol :: String -> Parser String
symbol = L.symbol sc


parens :: Parser a -> Parser a
parens = between (symbol "(") (symbol ")")


pNumber :: Parser HiExpr
pNumber = do
  sci <- L.signed sc (lexeme L.scientific)
  pure $ HiExprValue (HiValueNumber (realToFrac sci))


pBool :: Parser HiExpr
pBool = (symbol "true" >> pure (HiExprValue (HiValueBool True)))
    <|> (symbol "false" >> pure (HiExprValue (HiValueBool False)))


pNull :: Parser HiExpr
pNull = symbol "null" >> pure (HiExprValue HiValueNull)


pFunction :: Parser HiExpr
pFunction = do
  fun <- choice
    [ symbol "div" >> pure HiFunDiv
    , symbol "mul" >> pure HiFunMul
    , symbol "add" >> pure HiFunAdd
    , symbol "sub" >> pure HiFunSub
    , symbol "and" >> pure HiFunAnd
    , symbol "or" >> pure HiFunOr
    , symbol "less-than" >> pure HiFunLessThan
    , symbol "greater-than" >> pure HiFunGreaterThan
    , symbol "equals" >> pure HiFunEquals
    , symbol "not-less-than" >> pure HiFunNotLessThan
    , symbol "not-greater-than" >> pure HiFunNotGreaterThan
    , symbol "not-equals" >> pure HiFunNotEquals
    , symbol "not" >> pure HiFunNot
    , symbol "if" >> pure HiFunIf
    , symbol "length" >> pure HiFunLength
    , symbol "to-upper" >> pure HiFunToUpper
    , symbol "to-lower" >> pure HiFunToLower
    , symbol "reverse" >> pure HiFunReverse
    , symbol "trim" >> pure HiFunTrim
    , symbol "list" >> pure HiFunList
    , symbol "range" >> pure HiFunRange
    , symbol "fold" >> pure HiFunFold
    , symbol "pack-bytes" >> pure HiFunPackBytes
    , symbol "unpack-bytes" >> pure HiFunUnpackBytes
    , symbol "encode-utf8" >> pure HiFunEncodeUtf8
    , symbol "decode-utf8" >> pure HiFunDecodeUtf8
    , symbol "zip" >> pure HiFunZip
    , symbol "unzip" >> pure HiFunUnzip
    , symbol "serialise" >> pure HiFunSerialise
    , symbol "deserialise" >> pure HiFunDeserialise
    , symbol "read" >> pure HiFunRead
    , symbol "write" >> pure HiFunWrite
    , symbol "mkdir" >> pure HiFunMkDir
    , symbol "cd" >> pure HiFunChDir
    , symbol "parse-time" >> pure HiFunParseTime
    , symbol "rand" >> pure HiFunRand
    , symbol "echo" >> pure HiFunEcho
    , symbol "count" >> pure HiFunCount
    , symbol "keys" >> pure HiFunKeys
    , symbol "values" >> pure HiFunValues
    , symbol "invert" >> pure HiFunInvert
    ]
  pure (HiExprValue (HiValueFunction fun))

pString :: Parser HiExpr
pString = do
  str <- lexeme (char '"' *> manyTill L.charLiteral (char '"'))
  pure $ HiExprValue (HiValueString (T.pack str))

pList :: Parser HiExpr
pList = do
  _ <- symbol "["
  items <- sepBy exprParser (symbol ",")
  _ <- symbol "]"
  pure $ HiExprApply (HiExprValue (HiValueFunction HiFunList)) items

pDict :: Parser HiExpr
pDict = braces $ do
  entries <- sepBy dictEntry (symbol ",")
  pure (HiExprDict entries)
  where
    dictEntry = do
      key <- exprParser
      _ <- symbol ":"
      value <- exprParser
      pure (key, value)


braces :: Parser a -> Parser a
braces = between (symbol "{") (symbol "}")

identifier :: Parser String
identifier = lexeme ((:) <$> C.letterChar <*> many (C.alphaNumChar <|> char '-'))

pBytes :: Parser HiExpr
pBytes = do
  _ <- symbol "[#"
  sc
  bytes <- manyTill (try parseHexByte) (symbol "#]")
  pure $ HiExprValue (HiValueBytes (BS.pack bytes))
  where
    parseHexByte :: Parser Word8
    parseHexByte = do
      sc
      digits <- count 2 (satisfy isHexDigit)
      sc
      let hexToWord8 [d1, d2] =
            (hexDigitToInt d1 `shiftL` 4) .|. hexDigitToInt d2
          hexToWord8 _ = error "Invalid hex digit"
      return $ fromIntegral (hexToWord8 digits)

    hexDigitToInt :: Char -> Int
    hexDigitToInt c
      | c >= '0' && c <= '9' = ord c - ord '0'
      | c >= 'a' && c <= 'f' = ord c - ord 'a' + 10
      | c >= 'A' && c <= 'F' = ord c - ord 'A' + 10
      | otherwise = error "Invalid hex digit"

pCwd :: Parser HiExpr
pCwd = symbol "cwd" >> pure (HiExprValue (HiValueAction HiActionCwd))

pNow :: Parser HiExpr
pNow = symbol "now" >> pure (HiExprValue (HiValueAction HiActionNow))

atom :: Parser HiExpr
atom = pNumber <|> pBool <|> pNull <|> pString <|> pBytes <|> pList <|> pFunction <|> pNow <|> pCwd <|> pDict


operatorTable :: [[Expr.Operator Parser HiExpr]]
operatorTable =
  [ [ prefix "!" (HiExprValue (HiValueFunction HiFunNot)) ]
  , [ binaryL "*" (HiExprValue (HiValueFunction HiFunMul))
    , binaryL "/" (HiExprValue (HiValueFunction HiFunDiv))
    ]
  , [ binaryL "+" (HiExprValue (HiValueFunction HiFunAdd))
    , binaryL "-" (HiExprValue (HiValueFunction HiFunSub))
    ],
    [
      binaryN "==" (HiExprValue (HiValueFunction HiFunEquals))
    , binaryN "/=" (HiExprValue (HiValueFunction HiFunNotEquals))
    , binaryN "<=" (HiExprValue (HiValueFunction HiFunNotGreaterThan))
    , binaryN ">=" (HiExprValue (HiValueFunction HiFunNotLessThan))
    , binaryN "<"  (HiExprValue (HiValueFunction HiFunLessThan))
    , binaryN ">"  (HiExprValue (HiValueFunction HiFunGreaterThan))
    ]
  , [ binaryR "&&" (HiExprValue (HiValueFunction HiFunAnd)) ]
  , [ binaryR "||" (HiExprValue (HiValueFunction HiFunOr)) ]
  ]

opToken :: String -> Parser String
opToken "/"  = lexeme (try (string "/" <* notFollowedBy (char '=')))
opToken "<"  = lexeme (try (string "<" <* notFollowedBy (char '=')))
opToken ">"  = lexeme (try (string ">" <* notFollowedBy (char '=')))
opToken name = lexeme (try (string name))


binaryL :: String -> HiExpr -> Expr.Operator Parser HiExpr
binaryL name func = Expr.InfixL (makeInfix name func)

binaryR :: String -> HiExpr -> Expr.Operator Parser HiExpr
binaryR name func = Expr.InfixR (makeInfix name func)

binaryN :: String -> HiExpr -> Expr.Operator Parser HiExpr
binaryN name func = Expr.InfixN (makeInfix name func)

prefix :: String -> HiExpr -> Expr.Operator Parser HiExpr
prefix name func = Expr.Prefix (makePrefix name func)


makeInfix :: String -> HiExpr -> Parser (HiExpr -> HiExpr -> HiExpr)
makeInfix name func = do
  _ <- opToken name
  pure (\left right -> HiExprApply func [left, right])


makePrefix :: String -> HiExpr -> Parser (HiExpr -> HiExpr)
makePrefix name func = do
  _ <- (lexeme . string) name
  pure (\expr -> HiExprApply func [expr])

term :: Parser HiExpr
term = do
  base <- atom <|> parens exprParser
  postfixChain base

postfixChain :: HiExpr -> Parser HiExpr
postfixChain e =
      (do _ <- symbol "."
          fld <- identifier
          postfixChain (HiExprApply e [HiExprValue (HiValueString (T.pack fld))]))
  <|> (do args <- parens (sepBy exprParser (symbol ","))
          postfixChain (HiExprApply e args))
  <|> (do _ <- symbol "!"
          postfixChain (HiExprRun e))
  <|> pure e


exprParser :: Parser HiExpr
exprParser = Expr.makeExprParser term operatorTable




