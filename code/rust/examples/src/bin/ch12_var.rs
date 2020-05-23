#![no_std]
#![no_main]

#[panic_handler]
fn panic(_: &core::panic::PanicInfo) -> ! {
    loop {}
}

static mut G: i32 = 0xABCD;

#[no_mangle]
pub extern "C" fn get_g() -> i32 {
    unsafe { G }
}

#[no_mangle]
pub extern "C" fn set_g(g: i32) {
    unsafe {
        G = g;
    }
}

#[no_mangle]
pub extern "C" fn arr(i: usize) -> i32 {
    let a: [i32; 4] = [1, 2, 3, 4];
    if i < 4 {
        a[i]
    } else {
        0
    }
}
