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
}

pub struct S {
    a: i8,
    b: u8,
    c: i16,
    d: u16,
    e: i32,
    f: u32,
    g: i64,
    h: u64,
    i: f32,
    j: f64,
}

#[no_mangle]
pub extern "C" fn load_i32(s: &S) {
    unsafe {
        print_i32(s.a as i32);
        print_i32(s.b as i32);
        print_i32(s.c as i32);
        print_i32(s.d as i32);
        print_i32(s.e as i32);
        print_i32(s.f as i32);
    }
}

#[no_mangle]
pub extern "C" fn load_i64(s: &S) {
    unsafe {
        print_i64(s.a as i64);
        print_i64(s.b as i64);
        print_i64(s.c as i64);
        print_i64(s.d as i64);
        print_i64(s.e as i64);
        print_i64(s.f as i64);
        print_i64(s.g as i64);
        print_i64(s.h as i64);
    }
}

#[no_mangle]
pub extern "C" fn load_f(s: &S) {
    unsafe {
        print_f32(s.i);
        print_f64(s.j);
    }
}

#[no_mangle]
pub extern "C" fn store(s: &mut S, v: i64) {
    s.a = v as i8;
    s.b = v as u8;
    s.c = v as i16;
    s.d = v as u16;
    s.e = v as i32;
    s.f = v as u32;
    s.g = v as i64;
    s.h = v as u64;
    s.i = v as f32;
    s.j = v as f64;
}
