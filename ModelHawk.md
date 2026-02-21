# ModelHawk

## Intro

ModelHawk is a (proposed) standard way of connecting LLM-based tools (like Copilot, Claude, OpenCode) to tools that, for security purposes, monitor or control these LLM-based tools.

**Problem to solve:** We do not yet know how to gain assurance that LLM-based tools will not misbehave (due to, e.g., a prompt-injection attack).
Of course, there are many ideas for this, and more to come in the future. However, if we wanted to test a particular method of securing an LLM-based tool,
we don't have a standard way of hooking it up to such tools.

**Solution:** ModelHawk proposes a small protocol that LLM-based tools, and LLM security tools, can follow in order to connect to each other.

## Example

Suppose I run a team and I want to let my people use their favorite AI helpers (e.g. Claude CoWork) for their work. I want to prevent
these AI helpers from exfiltrating confidential calendar data (perhaps due to a prompt-injection attack).

So, I make (local network) service that implements the server part of ModelHawk, and I require my people to configure their AI helpers to connect to this security service as ModelHawk clients and notify this service about all HTTP tool uses (that is, each time a helper goes to a webpage).
My service can then do various things to check for exfiltration --- keyword search, asking another LLM, etc.

Now suppose that our security requirements get tighter, and I want all these AI helpers to ask permission before doing an HTTP request. This can be done by simply
changing how they connect to my monitoring service as ModelHawk clients, and of course modifying my service to return responses as appropriate.
