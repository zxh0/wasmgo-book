(module
  (import "env" "assert_eq_i32" (func $assert_eq_i32 (param i32 i32)))
  (start $main)
  (func $main (export "main")
    (call $assert_eq_i32 (call $sum (i32.const 1) (i32.const 100)) (i32.const 5050))
  )
  (func $sum (param $from i32) (param $to i32) (result i32)
    (local $n i32)

    (loop $l
      (; $n += $from ;)
      (local.set $n (i32.add (local.get $n) (local.get $from)))

      (; $from++ ;)
      (local.set $from (i32.add (local.get $from) (i32.const 1)))

      (; if $from <= $to { continue } ;)
      (br_if $l (i32.le_s (local.get $from) (local.get $to)))
    )

    (; return $n ;)
    local.get $n
  )
)
