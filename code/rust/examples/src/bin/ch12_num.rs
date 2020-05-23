#![no_std]
#![no_main]

#[panic_handler]
fn panic(_: &core::panic::PanicInfo) -> ! {
    loop {}
}

extern "C" {
    fn print_i32(n: i32);
    fn print_i64(n: i64);
    fn print_f32(n: f32);
    fn print_f64(n: f64);
    fn print_bool(n: bool);
}

#[no_mangle]
pub extern "C" fn num(a: i32, b: i32) {
    unsafe {
        print_i32(100);     // i32.const
        print_bool(a == 0); // i32.eqz
        print_bool(a >= b); // i32.ge_s
        print_i32(a * b);   // i32.mul
    }
}

#[no_mangle]
pub extern "C" fn conv(a: i32, b: i64, c: f32, d: f64) {
    unsafe {
        print_i32(b as i32); // i32.wrap_i64
        print_i32(c as i32); // i32.trunc_f32_s
        print_i64(a as i64); // i64.extend_i32_s
        print_f32(a as f32); // f32.convert_i32_s
        print_f32(d as f32); // f32.demote_f64
        print_f64(c as f64); // f64.promote_f32
        print_i32(f32::to_bits(c) as i32); // i32.reinterpret_f32
        print_i64(f64::to_bits(d) as i64); // i64.reinterpret_f64
        print_f32(f32::from_bits(a as u32)); // f32.reinterpret_i32
        print_f64(f64::from_bits(b as u64)); // f64.reinterpret_i64
    }
}
