// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/

// use tauri::tray::TrayIconBuilder;
use tauri::{
    tray::{MouseButtonState, TrayIconBuilder, TrayIconEvent},
    Manager,
};

use tauri::App;
use tauri_plugin_positioner::{Position, WindowExt};
use tauri_plugin_shell::ShellExt;

fn execute_polling(app: &mut App) {
    let sidecar_command = app.shell().sidecar("go-am-discord-rpc").unwrap();
    let (mut rx, mut _child) = sidecar_command.spawn().expect("Failed to spawn sidecar");
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_positioner::init())
        .plugin(tauri_plugin_shell::init())
        .setup(|app| {
            #[cfg(desktop)]
            {
                #[cfg(target_os = "macos")]
                let _ = app.set_activation_policy(tauri::ActivationPolicy::Accessory);

                let _tray = TrayIconBuilder::new()
                    .icon(app.default_window_icon().unwrap().clone())
                    .build(app)?;
                let _ = app.handle().plugin(tauri_plugin_positioner::init());
                tauri::tray::TrayIconBuilder::new()
                    .on_tray_icon_event(|tray_handle, event| {
                        tauri_plugin_positioner::on_tray_event(tray_handle.app_handle(), &event);
                    })
                    .build(app)?;
            }
            execute_polling(app);
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
            // Hide window on close, prevent
            tauri::WindowEvent::CloseRequested { api, .. } => {
                api.prevent_close();
                window.hide().unwrap();

                // Hide the dock icon when the window is closed (macOS-specific)
                // Hide the app from the dock (macOS specific)
                #[cfg(target_os = "macos")]
                let app = window.app_handle();
                let _ = app.set_activation_policy(tauri::ActivationPolicy::Accessory);
            }
            _ => {}
        })
        .plugin(tauri_plugin_opener::init())
        // .invoke_handler(tauri::generate_handler![greet, execute])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
