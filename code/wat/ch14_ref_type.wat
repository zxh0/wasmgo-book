;; --enable-reference-types
(module
  (table funcref (elem $f))
  (func $f
    (ref.is_null (ref.null func)) (drop)
    (ref.is_null (ref.func $f)) (drop)
  )
)