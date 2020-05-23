#![no_std]
#![no_main]

#[panic_handler]
fn panic(_: &core::panic::PanicInfo) -> ! {
    loop {}
}

extern "C" {
    fn random_i32() -> i32;
}

#[no_mangle]
pub extern "C" fn max(a: i32, b: i32) -> i32 {
    if a > b { a } else { b }
}

#[no_mangle]
pub extern "C" fn discard() {
  unsafe { let _ = random_i32(); }
}
