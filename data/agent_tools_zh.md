# Cursor Agent 工具完整总结

## 工具列表

| 序号 | 工具名称 | 参数消息 | 返回消息 | 主要参数 | 主要返回字段 |
|------|----------|----------|----------|----------|-------------|
| 1 | `READ_SEMSEARCH_FILES` | `ReadSemsearchFilesParams` | `ReadSemsearchFilesResult` | repository_info（仓库信息）、code_results[]（代码结果）、query（查询）、pr_references[]（PR引用）、pr_search_on（是否搜索PR） | code_results[]（代码结果）、all_files[]（全部文件）、missing_files[]（缺失文件）、knowledge_results[]（知识库结果）、pr_results[]（PR结果） |
| 2 | `RIPGREP_SEARCH` | `RipgrepSearchParams` | `RipgrepSearchResult` | options（搜索选项）、pattern_info（模式信息：正则、大小写、多行等） | internal（文件匹配[]、统计信息、退出码、消息） |
| 3 | `READ_FILE` | `ReadFileParams` | `ReadFileResult` | relative_workspace_path（相对路径）、read_entire_file（读取整个文件）、start_line_one_indexed（起始行）、end_line_one_indexed_inclusive（结束行）、max_lines（最大行数）、max_chars（最大字符数） | contents（内容）、outline（大纲）、start/end lines（行范围）、matching_cursor_rules[]（匹配规则）、file_git_context（Git上下文） |
| 4 | `LIST_DIR` | `ListDirParams` | `ListDirResult` | directory_path（目录路径） | files[]（文件列表：名称、是否目录、大小、最后修改时间、子项数、行数）、directory_relative_workspace_path（目录相对路径） |
| 5 | `EDIT_FILE` | `EditFileParams` | `EditFileResult` | relative_workspace_path（路径）、language（语言）、contents（内容）、old_string（旧字符串）、new_string（新字符串）、instructions（指令）、line_ranges[]（行范围）、blocking（阻塞）、allow_multiple_matches（允许多匹配） | diff（文件差异）、is_applied（已应用）、apply_failed（应用失败）、linter_errors[]（lint错误）、rejected（拒绝）、recoverable_error（可恢复错误）、human_review（人工审查） |
| 6 | `FILE_SEARCH` | `ToolCallFileSearchParams` | `ToolCallFileSearchResult` | query（查询关键词） | files[]（文件URI列表）、limit_hit（是否达上限）、num_results（结果数量） |
| 7 | `SEMANTIC_SEARCH_FULL` | `SemanticSearchFullParams` | `SemanticSearchFullResult` | repository_info（仓库信息）、query（查询）、include_pattern（包含模式）、exclude_pattern（排除模式）、top_k（前K个）、explanation（说明）、code_results[]（代码结果） | code_results[]（代码结果）、all_files[]（全部文件）、missing_files[]（缺失文件）、knowledge_results[]（知识库结果）、pr_results[]（PR结果） |
| 8 | `DELETE_FILE` | `DeleteFileParams` | `DeleteFileResult` | relative_workspace_path（相对路径） | rejected（拒绝）、file_non_existent（文件不存在）、file_deleted_successfully（删除成功） |
| 9 | `REAPPLY` | `ReapplyParams` | `ReapplyResult` | relative_workspace_path（相对路径） | diff（文件差异）、is_applied（已应用）、apply_failed（应用失败）、linter_errors[]（lint错误）、rejected（拒绝） |
| 10 | `RUN_TERMINAL_COMMAND_V2` | `RunTerminalCommandV2Params` | `RunTerminalCommandV2Result` | command（命令）、cwd（工作目录）、new_session（新会话）、options（超时、跳过AI检查等）、is_background（后台执行）、require_user_approval（需用户批准）、parsing_result（解析结果）、sandbox_policy（沙箱策略） | output（输出）、exit_code（退出码）、rejected（拒绝）、ended_reason（结束原因）、resulting_working_directory（结果目录）、output_location（输出位置） |
| 11 | `FETCH_RULES` | `FetchRulesParams` | `FetchRulesResult` | rule_names[]（规则名称列表） | rules[]（CursorRule：名称、描述、内容、always_apply） |
| 12 | `WEB_SEARCH` | `WebSearchParams` | `WebSearchResult` | search_term（搜索词） | references[]（引用：标题、URL、片段）、is_final（是否最终结果）、rejected（拒绝） |
| 13 | `MCP` | `MCPParams` | `MCPResult` | tools[]（工具列表：名称、描述、参数、服务器名）、file_output_threshold_bytes（文件输出阈值） | selected_tool（选中工具）、result（结果） |
| 14 | `SEARCH_SYMBOLS` | `SearchSymbolsParams` | `SearchSymbolsResult` | query（查询） | matches[]（匹配：名称、URI、范围、评分）、rejected（拒绝） |
| 15 | `BACKGROUND_COMPOSER_FOLLOWUP` | `BackgroundComposerFollowupParams` | `BackgroundComposerFollowupResult` | proposed_followup（提议的跟进）、bc_id（后台Composer ID） | proposed_followup（跟进内容）、is_sent（是否已发送） |
| 16 | `KNOWLEDGE_BASE` | `KnowledgeBaseParams` | `KnowledgeBaseResult` | knowledge_to_store（知识内容）、title（标题）、existing_knowledge_id（已有知识ID）、action（操作） | success（成功）、confirmation_message（确认消息）、id（ID） |
| 17 | `FETCH_PULL_REQUEST` | `FetchPullRequestParams` | `FetchPullRequestResult` | pull_number_or_commit_hash（PR号或提交哈希）、repo（仓库）、is_github（是否GitHub） | content（内容）、pr_number（PR号）、title（标题）、body（正文）、author（作者）、diff（差异）、comments[]（评论）、labels[]（标签）、state（状态） |
| 18 | `DEEP_SEARCH` | `DeepSearchParams` | `DeepSearchResult` | query（查询） | success（成功）、result（结果） |
| 19 | `CREATE_DIAGRAM` | `CreateDiagramParams` | `CreateDiagramResult` | content（图表内容） | error（错误，可选） |
| 20 | `FIX_LINTS` | `FixLintsParams` | `FixLintsResult` | （空参数） | file_results[]（文件结果：路径、差异、已应用、失败、错误、lint错误列表） |
| 21 | `READ_LINTS` | `ReadLintsParams` | `ReadLintsResult` | path（路径）、paths[]（路径列表） | path（路径）、linter_errors[]（lint错误）、linter_errors_by_file[]（按文件分组的lint错误） |
| 22 | `GO_TO_DEFINITION` | `GotodefParams` | `GotodefResult` | relative_workspace_path（路径）、symbol（符号）、start_line（起始行）、end_line（结束行） | definitions[]（定义列表：路径、完全限定名、符号类型、行范围、代码上下文行） |
| 23 | `TASK` | `TaskParams` | `TaskResult` | task_description（任务描述）、task_title（任务标题）、async（异步）、allowed_write_directories[]（允许写入目录）、model_override（模型覆盖） | completed_task_result（摘要、文件结果、用户中止）或 async_task_result（任务ID） |
| 24 | `AWAIT_TASK` | `AwaitTaskParams` | `AwaitTaskResult` | ids[]（任务ID列表） | task_results[]（任务ID、结果）、missing_task_ids[]（缺失任务ID） |
| 25 | `TODO_READ` | `TodoReadParams` | `TodoReadResult` | read（布尔值） | todos[]（待办项：内容、状态、ID、依赖） |
| 26 | `TODO_WRITE` | `TodoWriteParams` | `TodoWriteResult` | todos[]（待办项列表）、merge（合并） | success（成功）、ready_task_ids[]（就绪任务ID）、needs_in_progress_todos（需要进行中的待办）、final_todos[]（最终待办）、was_merge（是否合并） |
| 27 | `EDIT_FILE_V2` | `EditFileV2Params` | `EditFileV2Result` | relative_workspace_path（路径）、contents_after_edit（编辑后内容）、streaming_edit（流式编辑：文本/代码）、should_send_back_linter_errors（返回lint错误）、diff（差异） | contents_before_edit（编辑前内容）、file_was_created（文件已创建）、diff（差异）、rejected（拒绝）、linter_errors[]（lint错误）、result_for_model（模型结果） |
| 28 | `LIST_DIR_V2` | `ListDirV2Params` | `ListDirV2Result` | target_directory（目标目录）、ignore_globs[]（忽略规则）、should_enrich_terminal_metadata（终端元数据） | directory_tree_root（目录树节点：绝对路径、子目录列表、子文件列表、扩展名统计、文件数） |
| 29 | `READ_FILE_V2` | `ReadFileV2Params` | `ReadFileV2Result` | target_file（目标文件）、offset（偏移）、limit（限制）、chars_limit（字符限制）、effective_uri（有效URI）、enable_line_numbers（行号） | contents（内容）、num_characters_in_requested_range（字符数）、total_lines_in_file（总行数）、matching_cursor_rules[]（匹配规则）、images[]（图片） |
| 30 | `RIPGREP_RAW_SEARCH` | `RipgrepRawSearchParams` | `RipgrepRawSearchResult` | pattern（模式）、path（路径）、glob（通配符）、output_mode（输出模式）、context_before/after（上下文行）、case_insensitive（不区分大小写）、type（类型）、head_limit（头部限制）、ignore_globs[]（忽略规则） | success（模式、路径、输出模式、工作区结果{}、活动编辑器结果）或 error（错误） |
| 31 | `GLOB_FILE_SEARCH` | `GlobFileSearchParams` | `GlobFileSearchResult` | target_directory（目标目录）、glob_pattern（通配符模式） | directories[]（目录：绝对路径、文件列表、总文件数、是否截断） |
| 32 | `CREATE_PLAN` | `CreatePlanParams` | `CreatePlanResult` | plan（计划）、title（标题）、summary（摘要）、steps[]（步骤）、old_str/new_str（替换字符串）、todos[]（待办项）、overview（概述）、is_spec（是否规范）、phases[]（阶段） | plan_uri（计划URI）、accepted（接受，含最终待办）/ rejected（拒绝）/ modified（修改，含新计划和最终待办） |
| 33 | `LIST_MCP_RESOURCES` | `ListMcpResourcesParams` | `ListMcpResourcesResult` | server（服务器，可选） | resources[]（资源：URI、名称、描述、MIME类型、服务器、注解） |
| 34 | `READ_MCP_RESOURCE` | `ReadMcpResourceParams` | `ReadMcpResourceResult` | server（服务器）、uri（URI）、download_path（下载路径） | uri、name、description、mime_type、annotations{}、text 或 blob 内容 |
| 35 | `READ_PROJECT` | `ReadProjectParams` | `ReadProjectResult` | （空参数） | plan（计划内容） |
| 36 | `UPDATE_PROJECT` | `UpdateProjectParams` | `UpdateProjectResult` | string_replacements[]（字符串替换：旧值、新值）、summary（摘要） | success（成功）、updated_plan（更新后的计划） |
| 37 | `TASK_V2` | `TaskV2Params` | `TaskV2Result` | description（描述）、prompt（提示）、subagent_type（子代理类型）、model（模型）、name（名称） | agent_id（代理ID）、is_background（是否后台） |
| 38 | `CALL_MCP_TOOL` | `CallMcpToolParams` | `CallMcpToolResult` | server（服务器）、tool_name（工具名）、tool_args（工具参数，Struct类型） | server（服务器）、tool_name（工具名）、result（结果，Struct类型） |
| 39 | `APPLY_AGENT_DIFF` | `ApplyAgentDiffParams` | `ApplyAgentDiffResult` | agent_id（代理ID） | success（已应用变更列表）或 error（错误信息、已应用变更列表） |
| 40 | `ASK_QUESTION` | `AskQuestionParams` | `AskQuestionResult` | title（标题）、questions[]（问题：ID、提示、选项列表、允许多选）、run_async（异步执行） | answers[]（答案：问题ID、选中选项ID列表、自由文本）、is_async（是否异步） |
| 41 | `SWITCH_MODE` | `SwitchModeParams` | `SwitchModeResult` | from_mode_id（来源模式）、to_mode_id（目标模式）、explanation（说明） | from_mode_id、to_mode_id、auto_approved（自动批准）、user_approved（用户批准） |
| 42 | `GENERATE_IMAGE` | *（无参数消息）* | `GenerateImageResult` | *（特殊处理）* | success（文件路径、图片数据）或 error（错误信息） |
| 43 | `COMPUTER_USE` | `ComputerUseParams` | `ComputerUseResult` | actions[]（操作：鼠标移动、点击、按下/释放、拖拽、滚动、输入文字、按键、等待、截图、光标位置） | success（操作数量、耗时毫秒、截图、光标位置）或 error（错误信息） |
| 44 | `WRITE_SHELL_STDIN` | `WriteShellStdinArgs` | `WriteShellStdinResult` | shell_id（Shell ID）、chars（输入字符） | success（Shell ID、输入写入前的终端文件长度）或 error（错误信息） |
| 45 | `RECORD_SCREEN` | `RecordScreenArgs` | `RecordScreenResult` | mode（模式：开始/保存/丢弃录制）、tool_call_id、save_as_filename（保存文件名） | start_success（开始成功）/ save_success（路径、录制时长毫秒）/ discard_success（丢弃成功）/ failure（失败） |
| 46 | `WEB_FETCH` | `WebFetchParams` | `WebFetchResult` | url（网址） | url（网址）、markdown（Markdown内容）、error（错误） |
| 47 | `REPORT_BUGFIX_RESULTS` | `ReportBugfixResultsParams` | `ReportBugfixResultsResult` | summary（摘要）、results[]（结果：bug_id、bug_title、verdict判定、explanation说明） | success（结果列表）或 error（错误信息） |

---

## Agent 模式（来源：`AgentMode` 枚举）

| 值 | 模式 | 说明 |
|----|------|------|
| 0 | UNSPECIFIED | 未指定 |
| 1 | AGENT | 代理模式 |
| 2 | ASK | 询问模式 |
| 3 | PLAN | 计划模式 |
| 4 | DEBUG | 调试模式 |
| 5 | TRIAGE | 分诊模式 |
| 6 | PROJECT | 项目模式 |

---

## 子代理类型（来源：`SubagentType` 枚举）

| 值 | 类型 | 说明 |
|----|------|------|
| 0 | UNSPECIFIED | 未指定 |
| 1 | DEEP_SEARCH | 深度搜索 |
| 2 | FIX_LINTS | 修复Lint错误 |
| 3 | TASK | 任务 |
| 4 | SPEC | 规范 |

---

## 子代理参数与返回值

| 子代理类型 | 参数消息 | 返回消息 | 主要参数 | 主要返回字段 |
|-----------|---------|---------|---------|-------------|
| DEEP_SEARCH | `DeepSearchSubagentParams` | `DeepSearchSubagentReturnValue` | query（查询） | context_items[]（上下文项：文件、行范围、说明） |
| FIX_LINTS | `FixLintsSubagentParams` | `FixLintsSubagentReturnValue` | （空） | （空） |
| TASK | `TaskSubagentParams` | `TaskSubagentReturnValue` | task_description（任务描述）、allowed_write_directories[]（允许写入目录） | summary（摘要） |
| SPEC | `SpecSubagentParams` | `SpecSubagentReturnValue` | plan（计划） | summary（摘要）、string_replacements[]（字符串替换列表） |

---

## 关键支撑类型

### `AppliedAgentChange`（已应用的代理变更）

| 字段 | 类型 | 说明 |
|------|------|------|
| path | string | 文件路径 |
| change_type | 枚举 | CREATED（创建）/ MODIFIED（修改）/ DELETED（删除） |
| before_content | string（可选） | 变更前内容 |
| after_content | string（可选） | 变更后内容 |
| error | string（可选） | 错误信息 |
| message_for_model | string（可选） | 给模型的消息 |

### `SandboxPolicy`（沙箱策略）

| 字段 | 类型 | 说明 |
|------|------|------|
| type | 枚举 | INSECURE_NONE（无沙箱）/ WORKSPACE_READWRITE（工作区读写）/ WORKSPACE_READONLY（工作区只读） |
| network_access | bool（可选） | 是否允许网络访问 |
| additional_readwrite_paths | string[] | 额外可读写路径 |
| additional_readonly_paths | string[] | 额外只读路径 |
| block_git_writes | bool（可选） | 是否阻止Git写入 |
| disable_tmp_write | bool（可选） | 是否禁止临时目录写入 |
| allowlist_escalated | bool（可选） | 白名单是否已提升 |
| enable_shared_build_cache | bool（可选） | 是否启用共享构建缓存 |

### `ShellCommandParsingResult`（Shell命令解析结果）

| 字段 | 类型 | 说明 |
|------|------|------|
| parsing_failed | bool | 解析是否失败 |
| executable_commands | ExecutableCommand[] | 已解析的命令列表（名称、参数、完整文本） |
| has_redirects | bool | 是否有I/O重定向 |
| has_command_substitution | bool | 是否有命令替换 |

### `ComputerUseAction`（计算机操作，oneof类型）

| 操作 | 关键字段 | 说明 |
|------|---------|------|
| mouse_move | coordinate (x, y) | 鼠标移动 |
| click | coordinate, button, count, modifier_keys | 鼠标点击 |
| mouse_down | button | 鼠标按下 |
| mouse_up | button | 鼠标释放 |
| drag | path[] (坐标序列), button | 拖拽 |
| scroll | coordinate, direction, amount, modifier_keys | 滚动 |
| type | text | 输入文字 |
| key | key, hold_duration_ms | 按键 |
| wait | duration_ms | 等待 |
| screenshot | （空） | 截图 |
| cursor_position | （空） | 获取光标位置 |

### `TodoItem`（待办事项）

| 字段 | 类型 | 说明 |
|------|------|------|
| content | string | 待办内容 |
| status | string | 状态 |
| id | string | 唯一标识符 |
| dependencies | string[] | 依赖的其他待办ID |

### `BugfixResultItem`（Bug修复结果项）

| 字段 | 类型 | 说明 |
|------|------|------|
| bug_id | string | Bug标识符 |
| bug_title | string | Bug标题 |
| verdict | 枚举 | FIXED（已修复）/ FALSE_POSITIVE（误报）/ COULD_NOT_FIX（无法修复） |
| explanation | string | 说明 |

### `LinterError`（Lint错误）

| 字段 | 类型 | 说明 |
|------|------|------|
| message | string | 错误消息 |
| range | CursorRange | 错误范围（起始/结束行列） |
| source | string（可选） | 来源 |
| related_information | RelatedInformation[] | 相关信息 |
| severity | 枚举 | ERROR / WARNING / INFORMATION / HINT |
| is_stale | bool（可选） | 是否过期 |

### `EditFileResult.FileDiff`（文件差异）

| 字段 | 类型 | 说明 |
|------|------|------|
| chunks | ChunkDiff[] | 差异块列表（差异字符串、旧起始行、新起始行、旧行数、新行数、删除行数、新增行数） |
| editor | 枚举 | AI / HUMAN |
| hit_timeout | bool | 是否超时 |

### `ToolResultError`（工具结果错误）

| 字段 | 类型 | 说明 |
|------|------|------|
| client_visible_error_message | string | 客户端可见错误信息 |
| model_visible_error_message | string | 模型可见错误信息 |
| error_details | oneof | edit_file_error_details / search_replace_error_details |

### `ToolResultAttachments`（工具结果附件）

| 字段 | 类型 | 说明 |
|------|------|------|
| original_todos | TodoItem[] | 原始待办列表 |
| updated_todos | TodoItem[] | 更新后待办列表 |
| nudge_messages | NudgeMessage[] | 提醒消息 |
| should_show_todo_write_reminder | bool | 是否显示待办写入提醒 |
| todo_reminder_type | 枚举 | EVERY_10_TURNS / AFTER_EDIT |
| discovery_budget_reminder | DiscoveryBudgetReminder（可选） | 发现预算提醒（剩余轮数、发现力度） |

---

## 终端命令结束原因（`RunTerminalCommandEndedReason` 枚举）

| 值 | 原因 | 说明 |
|----|------|------|
| 0 | UNSPECIFIED | 未指定 |
| 1 | EXECUTION_COMPLETED | 执行完成 |
| 2 | EXECUTION_ABORTED | 执行中止 |
| 3 | EXECUTION_FAILED | 执行失败 |
| 4 | ERROR_OCCURRED_CHECKING_REASON | 检查原因时出错 |
| 5 | IDLE_TIMEOUT | 空闲超时 |

---

## 内置工具列表（`BuiltinTool` 枚举，旧版工具系统）

| 值 | 工具名 | 说明 |
|----|--------|------|
| 1 | SEARCH | 搜索 |
| 2 | READ_CHUNK | 读取代码块 |
| 3 | GOTODEF | 跳转定义 |
| 4 | EDIT | 编辑 |
| 5 | UNDO_EDIT | 撤销编辑 |
| 6 | END | 结束 |
| 7 | NEW_FILE | 新建文件 |
| 8 | ADD_TEST | 添加测试 |
| 9 | RUN_TEST | 运行测试 |
| 10 | DELETE_TEST | 删除测试 |
| 11 | SAVE_FILE | 保存文件 |
| 12 | GET_TESTS | 获取测试 |
| 13 | GET_SYMBOLS | 获取符号 |
| 14 | SEMANTIC_SEARCH | 语义搜索 |
| 15 | GET_PROJECT_STRUCTURE | 获取项目结构 |
| 16 | CREATE_RM_FILES | 创建/删除文件 |
| 17 | RUN_TERMINAL_COMMANDS | 运行终端命令 |
| 18 | NEW_EDIT | 新编辑 |
| 19 | READ_WITH_LINTER | 带Linter读取 |

---

## 交互更新消息类型（`InteractionUpdate` oneof）

| 消息类型 | 说明 |
|---------|------|
| TextDeltaUpdate | 文本增量更新 |
| PartialToolCallUpdate | 部分工具调用更新 |
| ToolCallDeltaUpdate | 工具调用增量更新 |
| ToolCallStartedUpdate | 工具调用开始 |
| ToolCallCompletedUpdate | 工具调用完成 |
| ThinkingDeltaUpdate | 思考增量更新 |
| ThinkingCompletedUpdate | 思考完成（含思考耗时毫秒） |
| UserMessageAppendedUpdate | 用户消息追加 |
| TokenDeltaUpdate | Token增量更新 |
| SummaryUpdate | 摘要更新 |
| SummaryStartedUpdate | 摘要开始 |
| SummaryCompletedUpdate | 摘要完成 |
| ShellOutputDeltaUpdate | Shell输出增量（stdout/stderr/exit/start） |
| HeartbeatUpdate | 心跳更新 |
| TurnEndedUpdate | 回合结束 |
| StepStartedUpdate | 步骤开始 |
| StepCompletedUpdate | 步骤完成（含步骤耗时毫秒） |

---

## 对话动作类型（`ConversationAction` oneof）

| 动作类型 | 说明 | 关键字段 |
|---------|------|---------|
| UserMessageAction | 用户消息 | user_message、request_context |
| ResumeAction | 恢复 | request_context |
| CancelAction | 取消 | reason（原因） |
| SummarizeAction | 总结 | （空） |
| ShellCommandAction | Shell命令 | shell_command、exec_id |
| StartPlanAction | 开始计划 | user_message、request_context、is_spec |
| ExecutePlanAction | 执行计划 | request_context、plan、execution_mode |
| AsyncAskQuestionCompletionAction | 异步提问完成 | original_tool_call_id、original_args、result |
| CancelSubagentAction | 取消子代理 | subagent_id |
