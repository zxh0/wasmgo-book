#!/usr/bin/env bash
set -ex

GO_DIR=$PWD/code/go
WAT_DIR=$PWD/code/wat
HW_DIR=$PWD/code/js

# compile .wat files
cd $WAT_DIR
rm -f *.wasm
for f in *.wat ; do
  if ! [[ $f =~ "ch14" ]] ; then 
    # echo "wat2wasm $f"
    wat2wasm "$f"
  else
    wat2wasm --enable-all "$f"
  fi
done

# test
cd $GO_DIR/ch01/wasm.go; go run wasm.go/cmd/wasmgo | grep "Hello, WebAssembly!"
cd $GO_DIR/ch02/wasm.go; go run wasm.go/cmd/wasmgo -d $HW_DIR/ch01_hw.wasm | grep "Version: 0x01"
cd $GO_DIR/ch03/wasm.go; go run wasm.go/cmd/wasmgo -d $HW_DIR/ch01_hw.wasm | grep "local.get 27"
cd $GO_DIR/ch05/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch05_param.wasm
cd $GO_DIR/ch05/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch05_num.wasm
cd $GO_DIR/ch05/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch05_cz.wasm
cd $GO_DIR/ch06/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch06_mem.wasm
cd $GO_DIR/ch07/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch07_global.wasm
cd $GO_DIR/ch07/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch07_local.wasm
cd $GO_DIR/ch07/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch07_max.wasm
cd $GO_DIR/ch07/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch07_sum.wasm
cd $GO_DIR/ch07/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch07_fib.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch07_sum.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch07_fib.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_eg4_add.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_eg5_calc.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_eg6_br.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_eg7_br_if.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_eg8_br_table.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_eg9_return.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_fac.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_cmp.wasm
cd $GO_DIR/ch08/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_sum.wasm
cd $GO_DIR/ch09/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch07_sum.wasm
cd $GO_DIR/ch09/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch07_fib.wasm
cd $GO_DIR/ch09/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_fac.wasm
cd $GO_DIR/ch09/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_cmp.wasm
cd $GO_DIR/ch09/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch08_sum.wasm
cd $GO_DIR/ch09/wasm.go; go run wasm.go/cmd/wasmgo $HW_DIR/ch01_hw.wasm | grep "Hello, World!"
cd $GO_DIR/ch09/wasm.go; go run wasm.go/cmd/wasmgo $WAT_DIR/ch09_calc.wasm
cd $GO_DIR/ch10/wasm.go; go run wasm.go/cmd/wasmgo $HW_DIR/ch01_hw.wasm | grep "Hello, World!"
cd $GO_DIR/ch11/wasm.go; go run wasm.go/cmd/wasmgo -c $HW_DIR/ch01_hw.wasm | grep "OK"

cd $GO_DIR/ch13/wasm.go
go run wasm.go/cmd/wasmgo $WAT_DIR/ch13_eg7_if.wasm

# hw
go run wasm.go/cmd/wasmgo -a $WAT_DIR/ch13_hw.wasm > hw.wasm.go
go build -buildmode=plugin -o hw.wasm.so hw.wasm.go
go run wasm.go/cmd/wasmgo hw.wasm.so | grep "Hello, World!"
rm hw.wasm.*

test_aot() {
  go run wasm.go/cmd/wasmgo -a $WAT_DIR/$1.wasm > $1.wasm.go
  go build -buildmode=plugin -o $1.wasm.so $1.wasm.go
  go run wasm.go/cmd/wasmgo $1.wasm.so
  rm $1.wasm.*
}
test_aot "ch05_param"
test_aot "ch05_num"
test_aot "ch06_mem"
test_aot "ch07_local"
test_aot "ch07_global"
test_aot "ch07_fib"
test_aot "ch07_max"
test_aot "ch07_sum"
test_aot "ch08_eg4_add"
test_aot "ch08_eg5_calc"
test_aot "ch08_eg6_br"
test_aot "ch08_eg7_br_if"
test_aot "ch08_eg8_br_table"
test_aot "ch08_eg9_return"
test_aot "ch08_cmp"
test_aot "ch08_sum"
test_aot "ch08_test"
test_aot "ch08_fac"
# test_aot "ch09_calc"


# remove .wasm files
cd $WAT_DIR
rm -f *.wasm

echo "OK!"
