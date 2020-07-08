;; --enable-reference-types
(module
  (table funcref (elem $f))
  (func $f
    (ref.is_null func (ref.null func)) (drop)
    (ref.is_null func (ref.func $f)) (drop)
  )
)