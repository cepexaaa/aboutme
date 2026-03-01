{-# LANGUAGE InstanceSigs #-}

module HW3.T3
  ( joinOption
  , joinExcept
  , joinAnnotated
  , joinList
  , joinFun
  ) where

import           HW3.T1 (Annotated (..), Except (..), Fun (..), List (..),
                         Option (..))

joinOption :: Option (Option a) -> Option a
joinOption None     = None
joinOption (Some x) = x

joinExcept :: Except e (Except e a) -> Except e a
joinExcept (Error e)   = Error e
joinExcept (Success x) = x

joinAnnotated :: Semigroup e => Annotated e (Annotated e a) -> Annotated e a
joinAnnotated ((x :# e1) :# e2) = x :# (e2 <> e1)

joinList :: List (List a) -> List a
joinList Nil = Nil
joinList (xs :. xss) = append xs (joinList xss)
  where
    append Nil ys        = ys
    append (x :. xs1) ys = x :. append xs1 ys

joinFun :: Fun i (Fun i a) -> Fun i a
joinFun (F f) = F (\i -> case f i of F g -> g i)
