#![no_std]
#![no_main]

#[panic_handler]
fn panic(_: &core::panic::PanicInfo) -> ! {
    loop {}
}

extern "C" {
    fn print_char(c: u8);
    fn print_i32(n: i32);
}

#[no_mangle]
pub extern "C" fn print_even(n: i32) {
    for x in 0..n {
        if x % 2 == 0 {
            unsafe { print_i32(x); }
        }
    }
}

#[no_mangle]
pub extern "C" fn print_ascii(n: i32) {
    unsafe {
        match n {
            0x61 => print_char('a' as u8),
            0x73 => print_char('s' as u8),
            0x6D => print_char('m' as u8),
            _    => print_char('?' as u8),
        }
    }
}
