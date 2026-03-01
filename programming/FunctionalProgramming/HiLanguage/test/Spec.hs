{-# LANGUAGE GeneralizedNewtypeDeriving #-}
{-# LANGUAGE OverloadedStrings          #-}

module Main (main) where

import           Control.Monad          (forM_, replicateM)
import           Data.Ratio             (denominator, numerator, (%))
import qualified Data.Sequence          as Seq
import qualified Data.Set               as Set
import qualified Data.Text              ()
import           HW5.Action             (HIO, HiPermission (..), runHIO)
import           HW5.Base
import           HW5.Evaluator
import           HW5.Parser
import           HW5.Pretty             (doc2String, prettyValue)
import           Test.Hspec

import           Control.Monad.Identity
import qualified Data.Map               as Map

newtype TestIdentity a = TestIdentity (Identity a)
  deriving (Functor, Applicative, Monad)

instance HiMonad TestIdentity where
  runAction _ = return HiValueNull

testEval :: HiExpr -> IO (Either HiError HiValue)
testEval expr =
  let hio = eval expr :: HIO (Either HiError HiValue)
      perms = Set.fromList [AllowRead, AllowWrite, AllowTime]
  in runHIO hio perms

isLeft :: Either a b -> Bool
isLeft (Left _) = True
isLeft _        = False

main :: IO ()
main = hspec $ do
  describe "Parser" $ do
    it "parses numbers" $ do
      parse "42" `shouldBe` Right (HiExprValue (HiValueNumber (42 % 1)))

    it "parses decimal fractions" $ do
      parse "3.14" `shouldBe` Right (HiExprValue (HiValueNumber (157 % 50)))

    it "parses negative numbers" $ do
      parse "-15" `shouldBe` Right (HiExprValue (HiValueNumber (-15 % 1)))

    it "parses function calls" $ do
      parse "add(1, 2)" `shouldBe`
        Right (HiExprApply (HiExprValue (HiValueFunction HiFunAdd))
               [HiExprValue (HiValueNumber (1 % 1)),
                HiExprValue (HiValueNumber (2 % 1))])

    it "parses nested function calls" $ do
      parse "div(add(10, 15.1), 3)" `shouldBe`
        Right (HiExprApply (HiExprValue (HiValueFunction HiFunDiv))
               [HiExprApply (HiExprValue (HiValueFunction HiFunAdd))
                 [HiExprValue (HiValueNumber (10 % 1)),
                  HiExprValue (HiValueNumber (151 % 10))],
                HiExprValue (HiValueNumber (3 % 1))])

  describe "Pretty printer" $ do
    it "renders integers" $ do
      let doc = doc2String $ prettyValue (HiValueNumber 42)
      doc `shouldBe` "42"

    it "renders decimal fractions" $ do
      let doc = doc2String $ prettyValue (HiValueNumber (157 % 50))
      doc `shouldBe` "3.14"

    it "renders simple fractions" $ do
      let doc = doc2String $ prettyValue (HiValueNumber (1 % 3))
      doc `shouldBe` "1/3"

    it "renders mixed fractions" $ do
      let doc = doc2String $ prettyValue (HiValueNumber (16 % 3))
      doc `shouldBe` "5 + 1/3"


  describe "Parser error cases" $ do
    it "fails on empty input" $ do
      parse "" `shouldSatisfy` isLeft

    it "fails on invalid syntax" $ do
      parse "add(1,,2)" `shouldSatisfy` isLeft
      parse "add(1,2" `shouldSatisfy` isLeft
      parse "add1,2)" `shouldSatisfy` isLeft

    it "fails on unknown identifiers" $ do
      parse "unknown(1, 2)" `shouldSatisfy` isLeft

    it "fails on invalid number format" $ do
      parse "12.34.56" `shouldSatisfy` isLeft



  describe "Infix operators" $ do
    it "parses 2 + 2" $ do
      parse "2 + 2" `shouldBe`
        Right (HiExprApply (HiExprValue (HiValueFunction HiFunAdd))
               [HiExprValue (HiValueNumber 2), HiExprValue (HiValueNumber 2)])

    it "parses <= operator" $ do
      let result = parse "2 <= 3"
      print result


    it "handles all infix operators" $ do
      let tests =
            [ ("2 / 3", HiValueNumber (2 % 3))
            , ("2 - 3", HiValueNumber (-1))
            , ("2 < 3", HiValueBool True)
            , ("2 > 3", HiValueBool False)
            , ("2 <= 3", HiValueBool True)
            , ("2 >= 3", HiValueBool False)
            , ("2 == 2", HiValueBool True)
            , ("2 /= 3", HiValueBool True)
            , ("true && false", HiValueBool False)
            , ("true || false", HiValueBool True)
            , ("!true", HiValueBool False)
            ]

      forM_ tests $ \(input, expected) -> do
        case parse input of
          Right expr -> do
            result <- testEval expr
            result `shouldBe` Right expected
          Left err -> expectationFailure $ "Should parse " ++ input ++ ": " ++ show err

  describe "File I/O" $ do
    it "parses read function" $ do
      case parse "read(\"test.txt\")" of
        Right expr -> expr `shouldBe`
          HiExprApply (HiExprValue (HiValueFunction HiFunRead))
            [HiExprValue (HiValueString "test.txt")]
        Left err -> error $ "Parse error: " ++ show err

    it "parses cwd value" $ do
      case parse "cwd" of
        Right expr -> expr `shouldBe`
          HiExprValue (HiValueAction HiActionCwd)
        Left err -> error $ "Parse error: " ++ show err

    it "parses run expression" $ do
      case parse "read(\"test.txt\")!" of
        Right expr -> expr `shouldBe`
          HiExprRun (HiExprApply (HiExprValue (HiValueFunction HiFunRead))
            [HiExprValue (HiValueString "test.txt")])
        Left err -> error $ "Parse error: " ++ show err


    it "parses now action" $ do
      case parse "now" of
        Right expr -> expr `shouldBe`
          HiExprValue (HiValueAction HiActionNow)
        Left err -> error $ "Parse error: " ++ show err

  describe "Random numbers" $ do
    it "parses rand function" $ do
      case parse "rand(0, 10)" of
        Right expr -> expr `shouldBe`
          HiExprApply (HiExprValue (HiValueFunction HiFunRand))
            [HiExprValue (HiValueNumber (0 % 1)),
            HiExprValue (HiValueNumber (10 % 1))]
        Left err -> error $ "Parse error: " ++ show err

    it "evaluates rand to action" $ do
      case parse "rand(0, 10)" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueAction (HiActionRand 0 10))
        Left err -> error $ "Parse error: " ++ show err

    it "rand produces random numbers in range" $ do

      let expr = HiExprRun (HiExprApply (HiExprValue (HiValueFunction HiFunRand))
                          [HiExprValue (HiValueNumber (0 % 1)),
                            HiExprValue (HiValueNumber (5 % 1))])


      results <- replicateM 100 $ do
        let resultIO = eval expr :: HIO (Either HiError HiValue)
        runHIO resultIO (Set.fromList [])


      forM_ results $ \result -> do
        case result of
          Right (HiValueNumber n) -> do
            let num = numerator n
            denominator n `shouldBe` 1
            num `shouldSatisfy` (\x -> x >= 0 && x <= 5)
          _ -> expectationFailure "Expected number in range 0-5"

    it "rand with invalid range returns error" $ do
      case parse "rand(10, 0)" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Left HiErrorInvalidArgument
        Left err -> error $ "Parse error: " ++ show err

  describe "Short-circuit evaluation and echo" $ do
    it "parses echo function" $ do
      case parse "echo(\"Hello\")" of
        Right expr -> expr `shouldBe`
          HiExprApply (HiExprValue (HiValueFunction HiFunEcho))
            [HiExprValue (HiValueString "Hello")]
        Left err -> error $ "Parse error: " ++ show err

    it "evaluates echo to action" $ do
      case parse "echo(\"Hello\")" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueAction (HiActionEcho "Hello"))
        Left err -> error $ "Parse error: " ++ show err

    it "short-circuit && with false" $ do
      case parse "false && echo(\"test\")!" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueBool False)
        Left err -> error $ "Parse error: " ++ show err

    it "short-circuit || with true" $ do
      case parse "true || echo(\"test\")!" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueBool True)
        Left err -> error $ "Parse error: " ++ show err

    it "if doesn't evaluate both branches" $ do
      case parse "if(true, \"OK\", echo(\"WTF\")!)" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueString "OK")
        Left err -> error $ "Parse error: " ++ show err

      case parse "if(false, echo(\"WTF\")!, \"OK\")" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueString "OK")
        Left err -> error $ "Parse error: " ++ show err

    it "string indexing with ||" $ do
      case parse "\"Hello\"(0) || \"Z\"" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueString "H")
        Left err -> error $ "Parse error: " ++ show err

      case parse "\"Hello\"(99) || \"Z\"" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueString "Z")
        Left err -> error $ "Parse error: " ++ show err

    it "bytes with &&" $ do
      case parse "[# 00 ff #] && echo(\"test\")" of
        Right expr -> do
          result <- (testEval expr)
          case result of
            Right (HiValueAction (HiActionEcho "test")) -> return ()
            _ -> expectationFailure $ "Expected echo action, got: " ++ show result
        Left err -> error $ "Parse error: " ++ show err

  describe "Dictionaries" $ do
    it "parses dictionary literal" $ do
      case parse "{ \"width\": 120, \"height\": 80 }" of
        Right expr -> expr `shouldBe`
          HiExprDict [(HiExprValue (HiValueString "width"),
                      HiExprValue (HiValueNumber (120 % 1))),
                      (HiExprValue (HiValueString "height"),
                      HiExprValue (HiValueNumber (80 % 1)))]
        Left err -> error $ "Parse error: " ++ show err

    it "accesses dictionary by key" $ do
      case parse "{ \"width\": 120, \"height\": 80 }(\"width\")" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueNumber (120 % 1))
        Left err -> error $ "Parse error: " ++ show err

    it "dot access" $ do
      case parse "{ \"width\": 120, \"height\": 80 }.width" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueNumber (120 % 1))
        Left err -> error $ "Parse error: " ++ show err

    it "count string" $ do
      case parse "count(\"XXXOX\")" of
        Right expr -> do
          result <- (testEval expr)
          let expected = HiValueDict (Map.fromList
                [(HiValueString "O", HiValueNumber (1 % 1)),
                (HiValueString "X", HiValueNumber (4 % 1))])
          result `shouldBe` Right expected
        Left err -> error $ "Parse error: " ++ show err

    it "keys and values" $ do
      case parse "keys({ \"width\": 120, \"height\": 80 })" of
        Right expr -> do
          result <- (testEval expr)
          let expected = HiValueList (Seq.fromList
                [HiValueString "height", HiValueString "width"])
          result `shouldBe` Right expected
        Left err -> error $ "Parse error: " ++ show err

      case parse "values({ \"width\": 120, \"height\": 80 })" of
        Right expr -> do
          result <- (testEval expr)
          let expected = HiValueList (Seq.fromList
                [HiValueNumber (80 % 1), HiValueNumber (120 % 1)])
          result `shouldBe` Right expected
        Left err -> error $ "Parse error: " ++ show err

    it "invert dictionary" $ do
      case parse "invert({ \"x\": 1, \"y\": 2, \"z\": 1 })" of
        Right expr -> do
          result <- (testEval expr)
          let expected = HiValueDict (Map.fromList
                [(HiValueNumber (1 % 1),
                  HiValueList (Seq.fromList [HiValueString "x", HiValueString "z"])),
                (HiValueNumber (2 % 1),
                  HiValueList (Seq.fromList [HiValueString "y"]))])
          result `shouldBe` Right expected
        Left err -> error $ "Parse error: " ++ show err

    it "count with dot access" $ do
      case parse "count(\"Hello World\").o" of
        Right expr -> do
          result <- (testEval expr)
          result `shouldBe` Right (HiValueNumber (2 % 1))
        Left err -> error $ "Parse error: " ++ show err
