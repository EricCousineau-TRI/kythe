/*
 * Copyright 2017 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"context"
	"flag"
	"fmt"

	cpb "kythe.io/kythe/proto/common_proto"
	xpb "kythe.io/kythe/proto/xref_proto"
)

type diagnosticsCommand struct {
	baseDecorCommand
}

func (diagnosticsCommand) Name() string     { return "diagnostics" }
func (diagnosticsCommand) Synopsis() string { return "list a file's diagnostics" }
func (diagnosticsCommand) Usage() string    { return "" }
func (c *diagnosticsCommand) SetFlags(flag *flag.FlagSet) {
	c.baseDecorCommand.SetFlags(flag)
}
func (c diagnosticsCommand) Run(ctx context.Context, flag *flag.FlagSet, api API) error {
	req, err := c.baseRequest(flag)
	if err != nil {
		return err
	}
	req.Diagnostics = true

	LogRequest(req)
	reply, err := api.XRefService.Decorations(ctx, req)
	if err != nil {
		return err
	}

	return c.displayDiagnostics(reply)
}

func (c diagnosticsCommand) displayDiagnostics(decor *xpb.DecorationsReply) error {
	if DisplayJSON {
		return PrintJSONMessage(decor)
	}

	for _, d := range decor.Diagnostic {
		fmt.Fprintf(out, "%s\t%s\n", spanToString(d.Span), d.Message)
	}
	return nil
}

func spanToString(s *cpb.Span) string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%d:%d-%d:%d",
		s.Start.GetLineNumber(), s.Start.GetColumnOffset(),
		s.End.GetLineNumber(), s.End.GetColumnOffset())
}
