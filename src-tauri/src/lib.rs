// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/

// use tauri::tray::TrayIconBuilder;
use tauri::{
    command,
    tray::{MouseButtonState, TrayIconBuilder, TrayIconEvent},
    AppHandle, Manager,
};

use reqwest::header::USER_AGENT;
use tauri_plugin_positioner::{Position, WindowExt};
use tauri_plugin_shell::process::CommandChild;
use tauri_plugin_shell::ShellExt;
use window_vibrancy::{apply_vibrancy, NSVisualEffectMaterial};

fn execute_polling(app: &AppHandle) -> CommandChild {
    let sidecar_command = app.shell().sidecar("go-am-discord-rpc").unwrap();
    let (mut _rx, child) = sidecar_command.spawn().expect("Failed to spawn sidecar");

    return child;
}

#[command]
async fn call_kill_api() -> Result<String, String> {
    // The URL of the kill endpoint (replace with your actual URL)
    let request_url = "http://localhost:8080/kill";

    let client = reqwest::Client::new();
    let response = client
        .get(request_url)
        .header(USER_AGENT, "rust-web-api-client")
        .send()
        .await
        .map_err(|e| e.to_string())?;

    if response.status().is_success() {
        Ok("Process killed successfully.".into())
    } else {
        // Return an error message as a String
        Err(format!(
            "Failed to call kill endpoint: {}",
            response.status()
        ))
    }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
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
                let _ = app.set_activation_policy(tauri::ActivationPolicy::Accessory);

                let _ = app.handle().plugin(tauri_plugin_positioner::init());
                let _ = TrayIconBuilder::new()
                    .icon(app.default_window_icon().unwrap().clone())
                    .on_tray_icon_event(|tray_handle, event| {
                        tauri_plugin_positioner::on_tray_event(tray_handle.app_handle(), &event);
                    })
                    .build(app)?;
            }

            let app_handle = app.app_handle();
            let window = app.get_webview_window("main").unwrap();
            let _child_proc = execute_polling(app_handle);

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
            _ => {}
        })
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_process::init())
        .invoke_handler(tauri::generate_handler![call_kill_api])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
