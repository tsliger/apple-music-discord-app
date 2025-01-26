// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/

// use tauri::tray::TrayIconBuilder;
use tauri::{
    menu::{Menu, MenuItem},
    tray::TrayIconBuilder,
    Manager,
};

#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .setup(|app| {
            let show_i = MenuItem::with_id(app, "show", "Show", true, None::<&str>)?;
            let quit_i = MenuItem::with_id(app, "quit", "Quit", true, None::<&str>)?;
            let menu = Menu::with_items(app, &[&show_i, &quit_i])?;
            let tray = TrayIconBuilder::new()
                .menu(&menu)
                .icon(app.default_window_icon().unwrap().clone())
                .build(app)?;
            Ok(())
        })
        .on_menu_event(|app, event| match event.id.as_ref() {
            "quit" => {
                println!("quit menu item was clicked");
                app.exit(0);
            }
            "show" => {
                // Show all windows and bring the app back to the dock
                if cfg!(target_os = "macos") {
                    app.show().unwrap();
                }
                let windows = app.webview_windows();
                for window in windows.values() {
                    window.show().unwrap();
                    window.set_focus().unwrap();
                }

                #[cfg(target_os = "macos")]
                app.set_activation_policy(tauri::ActivationPolicy::Regular);
            }
            _ => {
                println!("menu item {:?} not handled", event.id);
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
                app.set_activation_policy(tauri::ActivationPolicy::Accessory);
            }
            _ => {}
        })
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![greet])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
