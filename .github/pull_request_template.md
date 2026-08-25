| Status | Type | Related ticket(s) or issue(s)|
|---|---|---|
| CHOOSE ONE: ⏱️ WIP, ✅ Ready for review & deploy, 📖 Ready for review | CHOOSE ONE: 🔨 Refactor, 🐛 Bug, ✨ Feature, 🧨 Hotfix, ⚙️ Tooling | N/A |

<!--
> [!IMPORTANT]
> If needed, put here anything immediately relevant to reviewers, such as:
> - why it's WIP or not ready for deploy
> - why it needs to be reviewed ASAP
> - if there are any dependency PRs that can be leveraged
-->

---


## Description

### Overview
<!--
Focus on *why* this change is being made, which could include things like the business impact or general problems this solves.

If you have any relevant design documents like RFCs or diagrams, it would be good to link them here.
-->

### Out of scope
<!--
If applicable, this section should include work that is out of scope or deliberate omissions, ideally with links to follow up tickets.
-->

### Changelist
<!--
At a high level, list out the changes that were made to the codebase and, if relevant, external systems to make this PR work.
-->


## Testing
<!--
This is should be more than "go test ./... passes".

This should focus on things like:

- any unit tests that you added or modified
- if adding a new feature or fundamentally changing an existing one, any specification tests you added or modified
- any design choices you made that impacts tests

If you have to do ANY manual verification at all, please also put it here. Things like raw CLI output, screenshots of the AWS console or Postman, etc. are useful to document what we need to do to test and deploy so we can minimize tribal knowledge and hopefully figure out steps on how to automate this as much as possible in the future.
-->

- [ ] I followed TDD to my best ability (see [Canon TDD by Kent Beck](https://newsletter.kentbeck.com/p/canon-tdd))
- [ ] I verified that the specs in `specifications/` still match what the system actually does and no case has gone stale
