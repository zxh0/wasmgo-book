;; --enable-bulk-memory
;; --enable-reference-types
(module
  (func $f)

  (table $t1 1 8 funcref)
  (elem $e1 (offset (i32.const 1)) $f $f $f)

  (func $init (param $dst i32) (param $src i32) (param $size i32)
    (table.init $e1 (local.get $dst) (local.get $src) (local.get $size))
  )
  (func $copy (param $dst i32) (param $src i32) (param $size i32)
    (table.copy (local.get $dst) (local.get $src) (local.get $size))
  )
  (func $fill (param $dst i32) (param $size i32)
    (table.fill $t1 (local.get $dst) (ref.func $f) (local.get $size))
  )
  (func $drop
    (elem.drop $e1)
  )
)