use std::process::{Command, Child};
use std::sync::Mutex;
use tauri::Manager;

struct BackendChild(Mutex<Option<Child>>);

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
  tauri::Builder::default()
    .setup(|app| {
      if cfg!(debug_assertions) {
        app.handle().plugin(
          tauri_plugin_log::Builder::default()
            .level(log::LevelFilter::Info)
            .build(),
        )?;
      }

      // Spawns the Go backend on setup
      let child = if cfg!(debug_assertions) {
        Command::new("go")
          .args(["run", "main.go", "-port", "54321"])
          .current_dir("../src-go") // Svelte CLI runs from root, Tauri runs from src-tauri
          .spawn()
          .ok()
      } else {
        let exe_path = std::env::current_exe().unwrap();
        let exe_dir = exe_path.parent().unwrap();
        let backend_exe = exe_dir.join("winrice-backend.exe");
        Command::new(backend_exe)
          .args(["-port", "54321"])
          .spawn()
          .ok()
      };

      if let Some(c) = child {
        app.manage(BackendChild(Mutex::new(Some(c))));
      }

      Ok(())
    })
    .on_window_event(|window, event| {
      if let tauri::WindowEvent::Destroyed = event {
        // Kill Go backend on exit
        if let Some(state) = window.app_handle().try_state::<BackendChild>() {
          if let Ok(mut lock) = state.0.lock() {
            if let Some(mut child) = lock.take() {
              let _ = child.kill();
            }
          }
        }
      }
    })
    .run(tauri::generate_context!())
    .expect("error while running tauri application");
}

