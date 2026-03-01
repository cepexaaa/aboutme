{-# LANGUAGE OverloadedStrings #-}

module Main (main) where

import System.Console.Haskeline
import HW5.Parser (parse)
import HW5.Evaluator (eval)
import HW5.Pretty (prettyValue, prettyError, doc2String)
import qualified Data.Set as Set
import Text.Megaparsec (errorBundlePretty)
import Control.Monad.IO.Class (liftIO)
import HW5.Action (runHIO, HiPermission(..))

repl :: Set.Set HiPermission -> IO ()
repl perms = runInputT defaultSettings loop
  where
    loop :: InputT IO ()
    loop = do
      minput <- getInputLine "hi> "
      case minput of
        Nothing -> return ()
        Just ":q" -> return ()
        Just input -> do
          case parse input of
            Left err -> 
              outputStrLn $ "Parse error: " ++ (errorBundlePretty err)
            Right expr -> do
              result <- liftIO $ runHIO (eval expr) perms
              case result of
                Left err -> do
                  outputStrLn $ "Error: " ++ (doc2String $ prettyError err)
                Right val -> do
                  outputStrLn $ doc2String $ prettyValue val
          loop

main :: IO ()
main = do
  putStrLn "It's hi language!"
  putStrLn "Enter expressions to evaluate, or :q to quit."
  let perms = Set.fromList [AllowRead, AllowWrite, AllowTime]
  repl perms
