{-# LANGUAGE GeneralizedNewtypeDeriving #-}

module HW5.Action
( HIO(..)
, HiPermission(..)
, PermissionException(..)
) where

import           Control.Exception      (Exception, throwIO)
import           Control.Monad          (unless)
import           Control.Monad.IO.Class (MonadIO, liftIO)
import qualified Data.ByteString        as BS
import           Data.List              (sort)
import qualified Data.Sequence          as Seq
import           Data.Set               (Set, member)
import qualified Data.Text              as T
import qualified Data.Text.Encoding     as TE
import           Data.Time              (getCurrentTime)
import           HW5.Base
import           System.Directory       (createDirectory, doesFileExist,
                                         getCurrentDirectory, listDirectory,
                                         setCurrentDirectory)
import           System.Random          (randomRIO)


data HiPermission =
    AllowRead
  | AllowWrite
  | AllowTime
  deriving (Show, Eq, Ord)

data PermissionException =
  PermissionRequired HiPermission
  deriving (Show)


instance Exception PermissionException

newtype HIO a = HIO { runHIO :: Set HiPermission -> IO a }

instance Functor HIO where
  fmap f (HIO g) = HIO (\perms -> fmap f (g perms))

instance Applicative HIO where
  pure x = HIO (\_ -> pure x)
  (HIO f) <*> (HIO x) = HIO (\perms -> f perms <*> x perms)

instance Monad HIO where
  (HIO f) >>= k = HIO (\perms -> do
    a <- f perms
    let HIO g = k a
    g perms)

instance MonadIO HIO where
  liftIO io = HIO (\_ -> io)

instance HiMonad HIO where
  runAction action = HIO $ \perms -> do
    case action of
      HiActionCwd -> do
        unless (AllowRead `member` perms) $
          throwIO (PermissionRequired AllowRead)
        dir <- getCurrentDirectory
        return (HiValueString (T.pack dir))

      HiActionChDir path -> do
        unless (AllowRead `member` perms) $
          throwIO (PermissionRequired AllowRead)
        setCurrentDirectory path
        return HiValueNull

      HiActionRead path -> do
        unless (AllowRead `member` perms) $
          throwIO (PermissionRequired AllowRead)
        isFile <- doesFileExist path
        if isFile
          then do
            content <- BS.readFile path
            case TE.decodeUtf8' content of
              Right text -> return (HiValueString text)
              Left _     -> return (HiValueBytes content)
          else do
            entries <- listDirectory path
            let values = map (HiValueString . T.pack) (sort entries)
            return (HiValueList (Seq.fromList values))

      HiActionWrite path content -> do
        unless (AllowWrite `member` perms) $
          throwIO (PermissionRequired AllowWrite)
        BS.writeFile path content
        return HiValueNull

      HiActionMkDir path -> do
        unless (AllowWrite `member` perms) $
          throwIO (PermissionRequired AllowWrite)
        createDirectory path
        return HiValueNull

      HiActionNow -> do
        unless (AllowTime `member` perms) $
          throwIO (PermissionRequired AllowTime)
        time <- getCurrentTime
        return (HiValueTime time)

      HiActionRand low high -> do
        if low <= high
          then do
            randomValue <- randomRIO (low, high)
            return (HiValueNumber (fromIntegral randomValue))
          else
            throwIO (userError "Invalid range: low > high")

      HiActionEcho text -> do
        unless (AllowWrite `member` perms) $
          throwIO (PermissionRequired AllowWrite)
        putStrLn (T.unpack text)
        return HiValueNull
