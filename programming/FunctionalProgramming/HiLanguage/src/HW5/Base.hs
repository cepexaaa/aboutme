{-# LANGUAGE DeriveAnyClass     #-}
{-# LANGUAGE DeriveGeneric      #-}
{-# LANGUAGE StandaloneDeriving #-}

module HW5.Base
( HiFun(..)
, HiValue(..)
, HiExpr(..)
, HiError(..)
, HiAction(..)
, HiMonad(..)
) where

import           Codec.Serialise (Serialise)
import           Data.ByteString (ByteString)
import qualified Data.Foldable   as F
import qualified Data.Map        as Map
import           Data.Sequence   (Seq)
import           Data.Text       (Text)
import           Data.Time       (UTCTime)
import           GHC.Generics    (Generic)

data HiFun
  -- math
  = HiFunDiv
  | HiFunMul
  | HiFunAdd
  | HiFunSub
  -- bool
  | HiFunNot
  | HiFunAnd
  | HiFunOr
  | HiFunLessThan
  | HiFunGreaterThan
  | HiFunEquals
  | HiFunNotLessThan
  | HiFunNotGreaterThan
  | HiFunNotEquals
  | HiFunIf
  -- strings
  | HiFunLength
  | HiFunToUpper
  | HiFunToLower
  | HiFunReverse
  | HiFunTrim
  -- lists & folds
  | HiFunList
  | HiFunRange
  | HiFunFold
-- bytes & serialise
  | HiFunPackBytes
  | HiFunUnpackBytes
  | HiFunEncodeUtf8
  | HiFunDecodeUtf8
  | HiFunZip
  | HiFunUnzip
  | HiFunSerialise
  | HiFunDeserialise
-- read IO
  | HiFunRead
  | HiFunWrite
  | HiFunMkDir
  | HiFunChDir
-- time
  | HiFunParseTime
-- random
  | HiFunRand
-- circuit
  | HiFunEcho
-- dictionaries
  | HiFunCount
  | HiFunKeys
  | HiFunValues
  | HiFunInvert
  deriving (Show, Eq, Ord, Enum, Bounded, Generic, Serialise)

data HiValue
  = HiValueNumber Rational
  | HiValueBool Bool
  | HiValueNull
  | HiValueFunction HiFun
  | HiValueString Text
  | HiValueList (Seq HiValue)
  | HiValueBytes ByteString
  | HiValueAction HiAction
  | HiValueTime UTCTime
  | HiValueDict (Map.Map HiValue HiValue)
  deriving (Show, Eq, Generic, Serialise)

instance Ord HiValue where
  compare a b =
    case compare (tag a) (tag b) of
      EQ -> compareSame a b
      x  -> x
    where
      tag :: HiValue -> Int
      tag HiValueNull         = 0
      tag (HiValueBool _)     = 1
      tag (HiValueNumber _)   = 2
      tag (HiValueString _)   = 3
      tag (HiValueList _)     = 4
      tag (HiValueBytes _)    = 5
      tag (HiValueAction _)   = 6
      tag (HiValueTime _)     = 7
      tag (HiValueDict _)     = 8
      tag (HiValueFunction _) = 9

      compareSame :: HiValue -> HiValue -> Ordering
      compareSame HiValueNull HiValueNull = EQ
      compareSame (HiValueBool x) (HiValueBool y) = compare x y
      compareSame (HiValueNumber x) (HiValueNumber y) = compare x y
      compareSame (HiValueString x) (HiValueString y) = compare x y
      compareSame (HiValueList x) (HiValueList y) = compare (F.toList x) (F.toList y)
      compareSame (HiValueBytes x) (HiValueBytes y) = compare x y
      compareSame (HiValueAction x) (HiValueAction y) = compare x y
      compareSame (HiValueTime x) (HiValueTime y) = compare x y
      compareSame (HiValueDict x) (HiValueDict y) = compare (Map.toList x) (Map.toList y)
      compareSame (HiValueFunction x) (HiValueFunction y) = compare x y
      compareSame _ _ = EQ

data HiExpr
  = HiExprValue HiValue
  | HiExprApply HiExpr [HiExpr]
  | HiExprRun HiExpr
  | HiExprDict [(HiExpr, HiExpr)]
  deriving (Show, Eq)

data HiError
  = HiErrorInvalidArgument
  | HiErrorInvalidFunction
  | HiErrorArityMismatch
  | HiErrorDivideByZero
  deriving (Show, Eq)

data HiAction =
    HiActionRead  FilePath
  | HiActionWrite FilePath ByteString
  | HiActionMkDir FilePath
  | HiActionChDir FilePath
  | HiActionCwd
  | HiActionNow
  | HiActionRand Int Int
  | HiActionEcho Text
  deriving (Show, Eq, Generic, Serialise)

deriving instance Ord HiAction

class Monad m => HiMonad m where
  runAction :: HiAction -> m HiValue

