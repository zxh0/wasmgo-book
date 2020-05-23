#![no_std]
#![no_main]

#[panic_handler]
fn panic(_: &core::panic::PanicInfo) -> ! {
    loop {}
}

extern "C" {
    fn print_i64(n: i64);
}

#[no_mangle]
pub extern "C" fn main() {
    unsafe { print_i64(prime_count(5000000)) }
}

#[inline(never)]
#[no_mangle]
pub extern "C" fn prime_count(n: i64) -> i64 {
    let mut count = 0;
    for i in 0..n {
        if is_prime(i) {
            // unsafe { print_i64(i); }
            count += 1;
        }
    }
    return count;
}

#[inline(never)]
#[no_mangle]
pub extern "C" fn is_prime(n: i64) -> bool {
    if n <= 3 {
        return n > 1;
    } else if n % 2 == 0 || n % 3 == 0 {
        return false;
    }

    let mut i = 5;

    while i * i <= n {
        if i > 0 {
            if n % i == 0 || n % (i + 2) == 0 {
                return false;
            }
            i = i + 6;
        }
    }
    return true;
}
