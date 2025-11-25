# OCI 规范在 runx 项目中的实现详解

本文档详细说明 runx 项目如何实现和遵循 OCI (Open Container Initiative) 运行时规范。

---

## 📚 OCI 规范概述

OCI Runtime Specification 定义了容器运行时必须遵循的标准，包括：
- **状态定义**：容器生命周期中的标准状态
- **状态结构**：容器状态的 JSON 格式
- **生命周期操作**：创建、启动、停止、删除等

---

## 🔧 项目中的 OCI 实现

### 1. 核心依赖

项目使用 `github.com/opencontainers/runtime-spec/specs-go` 包，这是 OCI 规范的 Go 语言绑定。

```go
// 在多个文件中引入
import "github.com/opencontainers/runtime-spec/specs-go"
```

### 2. 状态结构的使用

#### 2.1 OCI State 结构嵌入

在 `pkg/sandbox/libcontainer/container_state.go` 中，项目将 OCI 的 `State` 结构嵌入到自定义的 `ContainerState` 中：

```go
type ContainerState struct {
    specs.State              // 嵌入 OCI State 结构
    Cmd     string           // 扩展字段：命令行
    Args    []string         // 扩展字段：参数
    Created time.Time        // 扩展字段：创建时间
}
```

**设计说明**：
- ✅ 通过嵌入（embedding）继承了 OCI 规范的所有标准字段
- ✅ 添加了额外的字段（Cmd, Args, Created）以满足项目需求
- ✅ 保持了与 OCI 规范的兼容性

#### 2.2 OCI State 字段说明

根据 OCI 规范，`specs.State` 包含以下标准字段：

```go
type State struct {
    Version     string            `json:"ociVersion"`  // OCI 规范版本
    ID          string            `json:"id"`          // 容器 ID
    Status      ContainerState    `json:"status"`      // 容器状态
    Pid         int               `json:"pid,omitempty"` // 进程 ID
    Bundle      string            `json:"bundle"`      // Bundle 路径
    Annotations map[string]string `json:"annotations,omitempty"` // 注解
}
```

---

## 🔄 状态映射机制

### 3.1 Linux 进程状态 → OCI 容器状态

项目实现了从 Linux 进程状态到 OCI 容器状态的映射机制，位于 `pkg/sandbox/system/proc.go`。

#### 步骤 1：读取 Linux 进程状态

```go
// 从 /proc/{pid}/stat 读取进程信息
func Stat(pid int) (stat Stat_t, err error) {
    bytes, err := ioutil.ReadFile(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
    // 解析进程状态字符: R (Running), S (Sleeping), Z (Zombie) 等
    return parseStat(string(bytes))
}
```

**Linux 进程状态**：
- `R` - Running（运行中）
- `S` - Sleeping（睡眠中）
- `Z` - Zombie（僵尸进程）
- `X` - Dead（已死亡）
- `T` - Stopped（已停止）

#### 步骤 2：状态转换逻辑

在 `GetContainerCmdAndStatus()` 函数中实现转换：

```go
func GetContainerCmdAndStatus(pid int) (string, specs.ContainerState, error) {
    var status = specs.StateRunning  // 默认状态
    cmd, err := os.Readlink("/proc/" + strconv.Itoa(pid) + "/exe")
    if err != nil {
        return "", specs.StateStopped, err  // 进程不存在 -> stopped
    }
    
    stat, err := Stat(pid)
    if err != nil {
        status = specs.StateStopped
    } else if stat.State == Zombie || stat.State == Dead {
        status = specs.StateStopped  // 僵尸或死亡 -> stopped
    }
    // 其他情况保持 running
    return cmd, status, err
}
```

**映射规则**：
```
Linux 进程状态          →  OCI 容器状态
─────────────────────────────────────────
Running (R)            →  StateRunning
Sleeping (S)           →  StateRunning (视为运行中)
Stopped (T)            →  StateStopped
Zombie (Z)             →  StateStopped
Dead (X)               →  StateStopped
进程不存在 / 读取失败    →  StateStopped
```

### 3.2 OCI 容器状态定义

根据 OCI 规范，容器状态有 4 种标准值：

```go
const (
    StateCreating ContainerState = "creating"  // 创建中
    StateCreated  ContainerState = "created"   // 已创建
    StateRunning  ContainerState = "running"   // 运行中
    StateStopped  ContainerState = "stopped"   // 已停止
)
```

**注意**：项目中主要使用了 `StateRunning` 和 `StateStopped` 两种状态。

---

## 💾 状态持久化

### 4.1 状态存储位置

容器状态被持久化到文件系统中：

```
/var/lib/wasm/WasmEdge/{pid}/config.json
```

### 4.2 状态数据结构

存储在 `config.json` 中的完整结构：

```go
type ContainerInfo struct {
    libcontainer.ContainerState  // 包含 OCI State + 扩展字段
    Labels      map[string]string `json:"labels,omitempty"`
    ContainerId string            `json:"container_id"`
}
```

### 4.3 状态写入示例

在 `wasm_edge.go` 的 `Init()` 方法中创建并写入状态：

```go
// 构建 OCI 标准状态
containerInfo := &types.ContainerInfo{
    ContainerState: libcontainer.ContainerState{
        State: specs.State{
            Version:     "1.0",              // OCI 版本
            Status:      status,             // OCI 状态（running/stopped）
            Pid:         pid,                // 进程 ID
            ID:          strconv.Itoa(pid),  // 容器 ID（使用 PID）
            Bundle:      w.Config.WASMFile,  // WASM 文件路径作为 Bundle
            Annotations: nil,                // 可选注解
        },
        Cmd: cmdString,                      // 扩展：命令行
    },
    ContainerId: strconv.Itoa(pid),         // 容器 ID
    Labels:      nil,                       // 可选标签
}

// 写入 JSON 文件
err = system.WriteContainerInfo(sandbox.WasmEdgeRuntimeRootPath, pid, containerInfo)
```

### 4.4 JSON 输出示例

`config.json` 的实际内容格式：

```json
{
  "ociVersion": "1.0",
  "id": "12345",
  "status": "running",
  "pid": 12345,
  "bundle": "/path/to/app.wasm",
  "annotations": null,
  "cmd": "/usr/bin/runx",
  "args": null,
  "created": "2024-01-01T00:00:00Z",
  "container_id": "12345",
  "labels": null
}
```

---

## 🔄 生命周期管理中的 OCI 实现

### 5.1 容器初始化（Init）

在 `WasmEdgeSandbox.Init()` 中：

```go
func (w *WasmEdgeSandbox) Init() (int, error) {
    // 1. 启动进程
    cmd := exec.Command(runx, args...)
    pid := cmd.Process.Pid
    
    // 2. 创建容器根目录
    system.GenerateContainerRootPath(sandbox.WasmEdgeRuntimeRootPath, pid)
    
    // 3. 获取初始状态并写入 OCI 格式的状态
    cmdString, status, err := system.GetContainerCmdAndStatus(pid)
    containerInfo := &types.ContainerInfo{
        ContainerState: libcontainer.ContainerState{
            State: specs.State{
                Version: "1.0",
                Status:  status,        // OCI 状态
                Pid:     pid,
                ID:      strconv.Itoa(pid),
                Bundle:  w.Config.WASMFile,
            },
        },
    }
    
    // 4. 持久化状态
    system.WriteContainerInfo(sandbox.WasmEdgeRuntimeRootPath, pid, containerInfo)
    
    return pid, nil
}
```

**对应 OCI 生命周期**：
- 相当于 OCI 的 `create` + `start` 操作
- 状态从 `creating` → `created` → `running`

### 5.2 状态查询（State）

在 `WasmEdgeSandbox.Sate()` 中：

```go
func (w *WasmEdgeSandbox) Sate() (*libcontainer.ContainerState, error) {
    // 1. 构建基础 OCI State
    state := specs.State{
        Version: "1.0",
        Status:  specs.StateRunning,  // 默认值
        Pid:     w.Config.Pid,
        ID:      strconv.Itoa(w.Config.Pid),
        Bundle:  "",
    }
    
    // 2. 从系统读取实时状态并更新
    cmd, status, err := system.GetContainerCmdAndStatus(w.Config.Pid)
    containerSate := &libcontainer.ContainerState{
        State: state,
    }
    containerSate.Cmd = cmd
    containerSate.State.Status = status  // 更新为实时 OCI 状态
    
    return containerSate, err
}
```

**特点**：
- ✅ 实时查询进程状态
- ✅ 返回 OCI 标准格式的状态
- ✅ 兼容 `runx wasm state` 命令的输出

### 5.3 容器列表（List）

通过读取文件系统目录实现：

```go
func (w *WasmEdgeSandbox) List() ([]string, error) {
    // 读取 /var/lib/wasm/WasmEdge/ 下的所有容器目录
    entries, err := ioutil.ReadDir(sandbox.WasmEdgeRuntimeRootPath)
    // 返回容器 ID 列表（即 PID 列表）
    return containers, nil
}
```

### 5.4 容器终止（Kill）

```go
func (w *WasmEdgeSandbox) Kill() error {
    // 1. 发送终止信号
    p, err := os.FindProcess(w.Config.Pid)
    p.Signal(syscall.SIGTERM)
    
    // 2. 清理容器根目录（包含状态文件）
    return os.RemoveAll(fmt.Sprintf("%s/%d", sandbox.WasmEdgeRuntimeRootPath, w.Config.Pid))
}
```

**对应 OCI 生命周期**：
- 相当于 OCI 的 `kill` + `delete` 操作
- 状态从 `running` → `stopped` → 删除

---

## 🎯 OCI 规范遵循情况

### ✅ 已实现的部分

1. **状态结构**：完整使用 `specs.State` 结构
2. **状态值**：使用标准状态值（running, stopped）
3. **状态持久化**：状态以 JSON 格式存储
4. **状态查询**：支持通过 PID 查询容器状态
5. **生命周期管理**：实现了 create、start、kill、delete 操作

### ⚠️ 部分实现的部分

1. **状态值**：
   - ✅ 使用了 `StateRunning` 和 `StateStopped`
   - ❌ 未使用 `StateCreating` 和 `StateCreated`
   - 💡 可以改进：在 Init() 时先设置 `StateCreating`，完成后设置 `StateCreated`，Start() 时设置为 `StateRunning`

2. **版本号**：
   - ⚠️ 硬编码为 `"1.0"`（代码中有 `// todo` 注释）
   - 💡 应该使用实际的 OCI 规范版本

3. **Bundle**：
   - ✅ 使用 WASM 文件路径作为 Bundle
   - ⚠️ 在查询状态时 Bundle 为空字符串

### ❌ 未实现的部分

1. **Annotations**：虽然字段存在，但始终为 `nil`
2. **完整生命周期**：未完全实现 `create` 和 `start` 的分离
3. **状态文件位置**：OCI 规范建议使用标准的运行时根目录格式

---

## 📊 实现架构图

```
┌─────────────────────────────────────────────────┐
│         OCI Runtime Specification               │
│    (github.com/opencontainers/runtime-spec)     │
└──────────────────┬──────────────────────────────┘
                   │
                   │ 嵌入和扩展
                   ▼
┌─────────────────────────────────────────────────┐
│     libcontainer.ContainerState                 │
│  ┌──────────────────────────────────────────┐  │
│  │  specs.State (OCI 标准状态)              │  │
│  │  - Version, ID, Status, Pid, Bundle      │  │
│  └──────────────────────────────────────────┘  │
│  + Cmd, Args, Created (扩展字段)                │
└──────────────────┬──────────────────────────────┘
                   │
                   │ 使用
        ┌──────────┴──────────┐
        │                     │
┌───────▼────────┐  ┌─────────▼────────┐
│  状态创建       │  │   状态查询       │
│  Init()        │  │   Sate()         │
│  - 写入文件    │  │   - 读取 /proc   │
│  - 持久化状态  │  │   - 实时状态映射 │
└────────────────┘  └──────────────────┘
        │                     │
        └──────────┬──────────┘
                   │
        ┌──────────▼──────────┐
        │  系统层状态映射      │
        │  GetContainerCmd    │
        │  AndStatus()        │
        │  Linux 进程状态     │
        │      ↓              │
        │  OCI 容器状态       │
        └─────────────────────┘
```

---

## 🔍 关键代码位置

| 功能 | 文件位置 | 关键函数/结构 |
|------|---------|--------------|
| OCI State 定义 | `vendor/github.com/opencontainers/runtime-spec/specs-go/state.go` | `State`, `ContainerState` |
| 状态结构扩展 | `pkg/sandbox/libcontainer/container_state.go` | `ContainerState` |
| 状态映射 | `pkg/sandbox/system/proc.go` | `GetContainerCmdAndStatus()` |
| 状态创建 | `pkg/sandbox/wasm/wasm_edge.go` | `Init()` |
| 状态查询 | `pkg/sandbox/wasm/wasm_edge.go` | `Sate()` |
| 状态持久化 | `pkg/sandbox/system/proc.go` | `WriteContainerInfo()` |
| 容器信息结构 | `pkg/types/types.go` | `ContainerInfo` |

---

## 💡 改进建议

### 1. 完整的状态生命周期

```go
// 在 Init() 中实现完整的状态转换
func (w *WasmEdgeSandbox) Init() (int, error) {
    // 1. 创建阶段
    containerInfo.State.Status = specs.StateCreating
    WriteContainerInfo(..., containerInfo)
    
    // 2. 初始化资源
    // ...
    
    // 3. 创建完成
    containerInfo.State.Status = specs.StateCreated
    WriteContainerInfo(..., containerInfo)
    
    // 4. 启动进程后
    containerInfo.State.Status = specs.StateRunning
    WriteContainerInfo(..., containerInfo)
    
    return pid, nil
}
```

### 2. 使用标准 OCI 版本

```go
const OciVersion = "1.1.0"  // 根据实际使用的规范版本

containerInfo.State.Version = OciVersion
```

### 3. 完善 Bundle 路径

```go
// 在 Sate() 中从 config.json 读取 Bundle，而不是留空
containerInfo := ReadContainerInfo(...)
containerSate.State.Bundle = containerInfo.State.Bundle
```

---

## 📚 参考资源

- [OCI Runtime Specification](https://github.com/opencontainers/runtime-spec)
- [OCI State JSON Schema](https://github.com/opencontainers/runtime-spec/blob/main/runtime.md#state)
- [opencontainers/runtime-spec Go Bindings](https://github.com/opencontainers/runtime-spec)

---

**总结**：runx 项目通过嵌入 OCI 规范的 `State` 结构、实现状态映射机制、持久化状态信息等方式，基本遵循了 OCI 运行时规范。虽然在完整生命周期和某些细节上还有改进空间，但核心的状态管理和持久化机制已经符合 OCI 标准。
