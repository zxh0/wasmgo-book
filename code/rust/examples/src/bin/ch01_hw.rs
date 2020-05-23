#![no_std]
#![no_main]

#[panic_handler]
fn panic(_: &core::panic::PanicInfo) -> ! {
    loop {}
}

extern "C" {
    fn print_char(c: u8);
}

#[no_mangle]
pub extern "C" fn main() {
    unsafe {
        let s = "Hello, World!\n";
        for c in s.as_bytes() {
            print_char(*c);
        }
    }
}
