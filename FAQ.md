# FAQ



Q：为什么我自己用`cargo`编译Rust例子产生的结果和书上写的不一样？

A：可能是因为`cargo`/`rustc`版本的问题。建议升级到最新版，然后修改Cargo.toml文件，关闭编译器优化后再试试：

```
[package]
name = "rust_examples"
version = "0.1.0"
edition = "2018"

[profile.release]
opt-level = 0 # 把这行注释打开
```



