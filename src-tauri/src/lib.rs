use tauri::{
    command,
    tray::{MouseButtonState, TrayIconBuilder, TrayIconEvent},
    AppHandle, Manager, State,
};

use std::net::TcpListener;
use std::process::Command;
use std::sync::Mutex;
use tauri_plugin_positioner::{Position, WindowExt};
use tauri_plugin_shell::ShellExt;
use tauri::Window;
use tauri_plugin_shell::process::CommandEvent;
use window_vibrancy::{apply_vibrancy, NSVisualEffectMaterial};

#[derive(Default)]
struct AppState {
    endpoint: String,
    go_child: u32,
}

fn kill_process(pid: &str) -> Result<(), Box<dyn std::error::Error>> {
    #[cfg(unix)]
    {
        println!("Sending INT signal to process with PID: {}", pid);

        let mut kill = Command::new("kill").args(["-s", "SIGINT", &pid]).spawn()?;
        kill.wait()?;
    }

    #[cfg(windows)]
    {
        println!("Sending taskkill to process with PID: {}", pid);

        let mut kill = Command::new("taskkill")
            .args(["/PID", &pid, "/F"])
            .spawn()?;
        kill.wait()?;
    }

    Ok(())
}

fn execute_polling(app: &AppHandle) {
    let open_port = find_open_port().unwrap();

    // let shell = app.shell();
    // let mut command = shell.command("go-am-discord-rpc").args(["42069"]);
    // command.args([open_port.clone()]); 

    let sidecar_command = app.shell().sidecar("go-am-discord-rpc").unwrap().args(["42069"]);
    let (mut rx, mut child) = sidecar_command.spawn().expect("Failed to spawn sidecar");

    // Send port number into std input
    let rest_endpoint = format!("http://localhost:{}/kill", open_port);

    println!("{}", rest_endpoint);

    // Set state
    let state = app.state::<Mutex<AppState>>();
    let mut state = state.lock().unwrap();

    state.endpoint = rest_endpoint;
    state.go_child = child.pid();

    // let window = app.get_webview_window("main").unwrap();
    println!("Port: {}", open_port);

    tauri::async_runtime::spawn(async move {
        while let Some(event) = rx.recv().await {
            if let CommandEvent::Stdout(line_bytes) = event {
                let line = String::from_utf8_lossy(&line_bytes);
                println!("{}", line);
                // window.emit("message", Some(format!("'{}'", line)))
                //     .expect("failed to emit event");
            }
        }
    });

    // println!("{}", child.pid().to_string());
}

fn find_open_port() -> Option<String> {
    let listener = TcpListener::bind("127.0.0.1:0").ok()?;
    let local_addr = listener.local_addr().ok()?;
    Some(local_addr.port().to_string())
}

#[command]
fn call_kill_api(state: State<'_, Mutex<AppState>>) {
    let state = state.lock().unwrap();

    let rest_endpoint = state.endpoint.to_string();

    println!("{}", rest_endpoint);

    let _response = reqwest::blocking::get(rest_endpoint);
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    let app = tauri::Builder::default()
        .plugin(tauri_plugin_autostart::init(
            tauri_plugin_autostart::MacosLauncher::LaunchAgent,
            Some(vec![]),
        ))
        .plugin(tauri_plugin_positioner::init())
        .plugin(tauri_plugin_shell::init())
        .setup(|app| {
            #[cfg(desktop)]
            {
                #[cfg(target_os = "macos")]
                app.set_activation_policy(tauri::ActivationPolicy::Accessory);

                let _ = app.handle().plugin(tauri_plugin_positioner::init());
                TrayIconBuilder::new()
                    .icon(app.default_window_icon().unwrap().clone())
                    .on_tray_icon_event(|tray_handle, event| {
                        tauri_plugin_positioner::on_tray_event(tray_handle.app_handle(), &event);
                    })
                    .build(app)?;
            }

            app.manage(Mutex::new(AppState::default()));

            let app_handle = app.app_handle();
            let window = app.get_webview_window("main").unwrap();
            execute_polling(app_handle);

            #[cfg(target_os = "macos")]
            apply_vibrancy(&window, NSVisualEffectMaterial::HudWindow, None, Some(16.0))
                .expect("Unsupported platform! 'apply_vibrancy' is only supported on macOS");

            Ok(())
        })
        .on_tray_icon_event(|app, event| {
            tauri_plugin_positioner::on_tray_event(app, &event);
            match event {
                TrayIconEvent::Click {
                    position: _,
                    button_state,
                    ..
                } => match button_state {
                    MouseButtonState::Down => {
                        if let Some(win) = app.get_webview_window("main") {
                            if win.is_visible().unwrap_or(false) {
                                let _ = win.hide();
                            } else {
                                let _ = win.move_window(Position::TrayCenter);
                                let _ = win.show();
                                let _ = win.set_focus();
                            }
                        }
                    }
                    MouseButtonState::Up => {}
                },
                _ => {}
            }
        })
        .on_window_event(|window, event| match event {
            tauri::WindowEvent::Focused(is_focused) => {
                if !is_focused {
                    window.hide().unwrap();
                }
            }
            tauri::WindowEvent::CloseRequested { .. } | tauri::WindowEvent::Destroyed { .. } => {
                let app = window.app_handle();
                let state = app.state::<Mutex<AppState>>();
                let state = state.lock().unwrap();

                let temp = state.go_child.to_string();
                println!("{}", temp);

                let _ = kill_process(temp.as_str());
            }
            _ => {}
        })
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_process::init())
        .invoke_handler(tauri::generate_handler![call_kill_api])
        .build(tauri::generate_context!())
        .expect("error while running tauri application");

    app.run(|app_handle, event| match event {
        tauri::RunEvent::ExitRequested { .. } => {
            let state = app_handle.state::<Mutex<AppState>>();
            let state = state.lock().unwrap();

            let temp = state.go_child.to_string();
            println!("{}", temp);

            let _ = kill_process(temp.as_str());
        }
        _ => {}
    });
}
