# Kiro CLI Research

This document analyzes Kiro CLI's interface, capabilities, and configuration options.

## Command Interface Analysis

### Main Commands and Capabilities

**Core Chat Commands:**
- `/quit` - Quit the application
- `/clear` - Clear conversation history
- `/agent` - Manage agents
- `/chat` - Manage saved conversations
- `/context` - Manage context files and view context window usage
- `/editor` - Open $EDITOR to compose prompts
- `/reply` - Open $EDITOR with most recent assistant message quoted
- `/compact` - Summarize conversation to free up context space
- `/tools` - View tools and permissions
- `/model` - Select a model for current conversation session

**Specialized Features:**
- `/code` - Code intelligence with LSP integration
- `/plan` - Switch to Plan agent for breaking down ideas
- `/todos` - View, manage, and resume to-do lists
- `/knowledge` - Store and retrieve information in knowledge base (experimental)
- `/experiment` - Toggle experimental features
- `/tangent` - Conversation checkpointing for exploring side topics

### Session Management

**Starting Sessions:**
- Install: `curl -fsSL https://cli.kiro.dev/install | bash`
- Start: `kiro-cli chat` (starts interactive session)
- With agent: `kiro-cli chat --agent myagent`

**Session Features:**
- Automatic conversation persistence
- Context management with `/context add <path>`
- Model selection with `/model`
- Agent switching with `/agent`

### Configuration System

**Settings Management:**
- Command: `kiro-cli settings` (from terminal, NOT in chat)
- View all: `kiro-cli settings`
- Set value: `kiro-cli settings key value`
- Boolean example: `kiro-cli settings chat.enableThinking true`

**Key Configuration Options:**
- `chat.defaultAgent` - Default agent configuration
- `chat.defaultModel` - Default AI model for conversations
- `chat.enableThinking` - Enable thinking tool for complex reasoning
- `chat.enableKnowledge` - Enable knowledge base functionality
- `chat.enableCodeIntelligence` - Enable code intelligence with LSP
- `chat.enableTangentMode` - Enable tangent mode feature
- `api.timeout` - API request timeout in seconds
- `telemetry.enabled` - Enable/disable telemetry collection

## Built-in Tools

**File Operations:**
- `fs_read` - Read files, directories, and images
- `fs_write` - Create and edit files
- `execute_bash` - Execute shell commands

**Development Tools:**
- `code` - Code intelligence with LSP support
- `use_aws` - Make AWS CLI API calls
- `grep` and `glob` - Search and pattern matching

**Productivity Tools:**
- `knowledge` - Store/retrieve information across sessions
- `todo_list` - Create and manage TODO lists
- `thinking` - Internal reasoning mechanism
- `use_subagent` - Delegate tasks to specialized agents

## Agent System

**Agent Locations:**
- Local (workspace): `.kiro/agents/`
- Global (user-wide): `~/.kiro/agents/`
- Precedence: Local agents override global agents

**Agent Configuration:**
- JSON format with tools, resources, and settings
- Tool permissions and restrictions
- Custom prompts and behavior

## Experimental Features

**Available Experiments:**
- **Checkpointing** - Session-scoped file change tracking
- **Context Usage Percentage** - Show context window usage in prompt
- **Knowledge** - Persistent context storage across sessions
- **Thinking** - Complex reasoning with step-by-step processes
- **Tangent Mode** - Conversation checkpointing for side topics
- **Delegate** - Background task management with subagents
- **TODO Lists** - Persistent task management

**Enable via:** `/experiment` command in chat
