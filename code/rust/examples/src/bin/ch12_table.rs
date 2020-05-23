#![no_std]
#![no_main]

#[panic_handler]
fn panic(_: &core::panic::PanicInfo) -> ! {
    loop {}
}

type Binop = fn(f32, f32) -> f32;
fn add(a: f32, b: f32) -> f32 { a + b }
fn sub(a: f32, b: f32) -> f32 { a - b }
fn mul(a: f32, b: f32) -> f32 { a * b }
fn div(a: f32, b: f32) -> f32 { a / b }

#[no_mangle]
pub extern "C" fn calc(op: usize, a: f32, b: f32) -> f32 {
    get_fn(op)(a, b)
}

fn get_fn(op: usize) -> Binop {
    match op {
        1 => add,
        2 => sub,
        3 => mul,
        _ => div,
    }
}
