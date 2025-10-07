# Benchy

Benchy is a modern, cross-platform system benchmarking tool written in Go with a sleek Qt-based interface. It measures CPU, memory, disk, and system performance through real-time, visually rich tests—presented in an intuitive Geekbench-style results dialog.

# Screenshots
<img width="800" alt="benchy_running_tests" src="https://github.com/user-attachments/assets/5c5b29be-1e91-4f57-9de4-6ec7226cc2c8" />
<img width="800" alt="benchy_tests_completed" src="https://github.com/user-attachments/assets/da9d7d20-acf3-406d-89b7-9892cb79a8e1" />
<img width="800" alt="benchy_results" src="https://github.com/user-attachments/assets/17758b64-837a-48b3-b902-0d89815e7b85" />


# Features

- Single-Core and Multi-Core tabs
- Per-test **sub-scores** and a big **Overall** tile (geometric mean, baseline 2500)
- A simple **bar chart** of sub-scores per tab
- Export JSON for each tab
- Dark mode palette (auto-applied)
- Section tiles (CPU / Memory / Storage / Image)
- Improved in-app bar chart:
  - Axis ticks and labels
  - Mouse hover with value tooltip

## Build
```bash
make
```
