// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package commands

import (
	"strings"
	stdtesting "testing"

	"github.com/juju/tc"

	"github.com/juju/juju/internal/testing"
)

type HelpHookCommandsSuite struct {
	testing.FakeJujuXDGDataHomeSuite
}

func TestHelpHookCommandsSuite(t *stdtesting.T) {
	tc.Run(t, &HelpHookCommandsSuite{})
}

func (suite *HelpHookCommandsSuite) SetUpTest(c *tc.C) {
	suite.FakeJujuXDGDataHomeSuite.SetUpTest(c)
	setFeatureFlags("")
}

func (suite *HelpHookCommandsSuite) TestHelpHookCommandsHelp(c *tc.C) {
	output := badrun(c, 0, "help", "help-hook-commands")
	c.Assert(output, tc.Equals, `Usage: juju help-hook-commands [hook]

Summary:
Show help on a Juju charm hook command.

Global Options:
--debug  (= false)
    Equivalent to --show-log --logging-config=<root>=DEBUG
-h, --help  (= false)
    Show help on a command or other topic.
--logging-config (= "")
    Specify log levels for modules
--quiet  (= false)
    Show no informational output
--show-log  (= false)
    If set, write the log file to stderr
--verbose  (= false)
    Show more verbose output

Details:
Juju charms have access to a set of built-in helpers known as 'hook-commands,'
which allow them to inspect their runtime environment.
The currently available charm hook commands include:

    application-version-set  Specifies which version of the application is deployed.
    close-port               Registers a request to close a port or port range.
    config-get               Prints application configuration.
    credential-get           Accesses cloud credentials.
    goal-state               Prints the status of the charm's peers and related units.
    is-leader                Prints application leadership status.
    juju-log                 Writes a message to Juju logs.
    juju-reboot              Reboots the host machine.
    network-get              Gets network config.
    open-port                Registers a request to open a port or port range.
    opened-ports             Lists all ports or port ranges opened by the unit.
    relation-get             Get relation settings.
    relation-ids             Lists all relation IDs for the given endpoint.
    relation-list            Lists relation units.
    relation-model-get       Gets details about the model housing a related application.
    relation-set             Sets relation settings.
    resource-get             Gets the path to the locally cached resource file.
    secret-add               Adds a new secret.
    secret-get               Gets the content of a secret.
    secret-grant             Grants access to a secret.
    secret-ids               Prints secret IDs.
    secret-info-get          Gets a secret's metadata info.
    secret-remove            Removes an existing secret.
    secret-revoke            Revokes access to a secret.
    secret-set               Updates an existing secret.
    state-delete             Deletes server-side-state key-value pairs.
    state-get                Prints server-side-state value.
    state-set                Sets server-side-state values.
    status-get               Prints status information.
    status-set               Sets status information.
    storage-add              Adds storage instances.
    storage-get              Prints information for the storage instance with the specified ID.
    storage-list             Lists storage attached to the unit.
    unit-get                 Prints public-address or private-address.

Examples:

For help on a specific hook command, supply the name of that hook command, for example:

    juju help-hook-commands unit-get

See also:
 - help
 - help-action-commands
`)
}

var expectedCommands = []string{
	"application-version-set",
	"close-port",
	"config-get",
	"credential-get",
	"goal-state",
	"is-leader",
	"juju-log",
	"juju-reboot",
	"network-get",
	"open-port",
	"opened-ports",
	"relation-get",
	"relation-ids",
	"relation-list",
	"relation-model-get",
	"relation-set",
	"resource-get",
	"secret-add",
	"secret-get",
	"secret-grant",
	"secret-ids",
	"secret-info-get",
	"secret-remove",
	"secret-revoke",
	"secret-set",
	"state-delete",
	"state-get",
	"state-set",
	"status-get",
	"status-set",
	"storage-add",
	"storage-get",
	"storage-list",
	"unit-get",
}

func (suite *HelpHookCommandsSuite) TestHelpHookCommands(c *tc.C) {
	output := badrun(c, 0, "help-hook-commands")
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for i, line := range lines {
		command := strings.Fields(line)[0]
		lines[i] = command
	}
	c.Assert(lines, tc.DeepEquals, expectedCommands)
}

func (suite *HelpHookCommandsSuite) TestHelpHookCommandsName(c *tc.C) {
	output := badrun(c, 0, "help-hook-commands", "relation-get")
	expectedHelp := `Usage: relation-get \[options\] <key> <unit id>

Summary:
Get relation settings.

Options:
(.|\n)*

Details:
` + "`relation-get`" + ` prints the value(.|\n)*`
	c.Assert(output, tc.Matches, expectedHelp)
}
