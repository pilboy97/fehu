package main

import "cli"

var parser *cli.Parser

var rootState = &cli.State{}
var openState = &cli.State{
	Action: true,
	Name:   "open",
	Manual: "open a fehu file",
	Pat:    "^open$",
	Param:  1,
}
var closeState = &cli.State{
	Action: true,
	Name:   "close",
	Manual: "close a fehu file",
	Pat:    "^close$",
}
var quitState = &cli.State{
	Action: true,
	Name:   "quit",
	Manual: "exit fehu system",
	Pat:    "^quit$",
}
var calcState = &cli.State{
	Action: true,
	Name:   "calc",
	Manual: "calculate a statement",
	Pat:    "^calc$",
	Param:  1,
}
var defState = &cli.State{
	Action: true,
	Name:   "def",
	Manual: "define variable",
	Pat:    "^def$",
	Param:  2,
}
var newState = &cli.State{
	Name:   "new",
	Manual: "create something",
	Pat:    "^new$",
}
var getState = &cli.State{
	Name:   "get",
	Manual: "find something",
	Pat:    "^get$",
}
var altState = &cli.State{
	Name:   "alt",
	Manual: "alter something",
	Pat:    "^alt$",
}
var delState = &cli.State{
	Name:   "del",
	Manual: "delete something",
	Pat:    "^del$",
}

var newAccState = &cli.State{
	Action: true,
	Name:   "new acc",
	Manual: "create a new account",
	Pat:    "^acc$",
	Param:  1,
	Flags: []*cli.Flag{
		{
			Name:    "desc",
			Manual:  "set description",
			NamePat: "^(d|desc)$",
			ValPat:  "^.*$",
		},
	},
}
var getAccState = &cli.State{
	Action: true,
	Name:   "get acc",
	Manual: "find some account",
	Pat:    "^acc$",
}
var getAccByNameState = &cli.State{
	Action: true,
	Name:   "get acc name",
	Manual: "find an account by name",
	Pat:    "^name$",
	Param:  1,
}
var getAccChildState = &cli.State{
	Action: true,
	Name:   "get acc child",
	Manual: "find childs of acc",
	Pat:    "^child$",
	Param:  1,
}
var getAccByDescState = &cli.State{
	Action: true,
	Name:   "get acc desc",
	Manual: "find an account by desc",
	Pat:    "^desc$",
	Param:  1,
}
var altAccState = &cli.State{
	Action: true,
	Name:   "alt acc",
	Manual: "alter an account",
	Pat:    "^acc$",
	Param:  1,
	Flags: []*cli.Flag{
		{
			Name:    "desc",
			Manual:  "set description",
			NamePat: "^(d|desc)$",
			ValPat:  "^.*$",
		},
	},
}
var altAccRenameState = &cli.State{
	Action: true,
	Name:   "alt acc rename",
	Manual: "rename an account",
	Pat:    "^rename$",
	Param:  2,
}
var delAccState = &cli.State{
	Action: true,
	Name:   "del acc",
	Manual: "delete an account",
	Pat:    "^acc$",
	Param:  1,
}

var newTxnState = &cli.State{
	Action: true,
	Name:   "new txn",
	Manual: "create a new transaction",
	Pat:    "^txn$",
	Param:  1,
	Flags: []*cli.Flag{
		{
			Name:    "desc",
			Manual:  "set description",
			NamePat: "^(d|desc)$",
			ValPat:  "^.*$",
		},
		{
			Name:    "time",
			Manual:  "set timestamp",
			NamePat: "^(t|time)$",
			ValPat:  "^\\d{4}-\\d{2}-\\d{2};\\d{2}:\\d{2}:\\d{2}$",
		},
	},
}
var getTxnState = &cli.State{
	Action: true,
	Name:   "get txn",
	Manual: "find some transaction",
	Pat:    "^txn$",
	Flags: []*cli.Flag{
		{
			Name:    "save",
			Manual:  "save table",
			NamePat: "^(s|save)$",
			ValPat:  `^~?(\p{L}(\p{L}|\d)*:)*(\p{L}(\p{L}|\d)*)$`,
		},
	},
}
var getTxnByIDState = &cli.State{
	Action: true,
	Name:   "get txn id",
	Manual: "find a transaction by ID",
	Pat:    "^id$",
	Param:  1,
	Flags: []*cli.Flag{
		{
			Name:    "save",
			Manual:  "save table",
			NamePat: "^(s|save)$",
			ValPat:  `^~?(\p{L}(\p{L}|\d)*:)*(\p{L}(\p{L}|\d)*)$`,
		},
	},
}
var getTxnByTimeState = &cli.State{
	Action: true,
	Name:   "get txn time",
	Manual: "find an transaction within period",
	Pat:    "^time$",
	Param:  1,
	Flags: []*cli.Flag{
		{
			Name:    "save",
			Manual:  "save table",
			NamePat: "^(s|save)$",
			ValPat:  `^~?(\p{L}(\p{L}|\d)*:)*(\p{L}(\p{L}|\d)*)$`,
		},
	},
}
var getTxnByDescState = &cli.State{
	Action: true,
	Name:   "get txn desc",
	Manual: "find an transaction by desc",
	Pat:    "^desc$",
	Param:  1,
	Flags: []*cli.Flag{
		{
			Name:    "save",
			Manual:  "save table",
			NamePat: "^(s|save)$",
			ValPat:  `^~?(\p{L}(\p{L}|\d)*:)*(\p{L}(\p{L}|\d)*)$`,
		},
	},
}
var altTxnState = &cli.State{
	Action: true,
	Name:   "alt txn",
	Manual: "alter some transaction",
	Pat:    "^txn$",
	Param:  1,
	Flags: []*cli.Flag{
		{
			Name:    "desc",
			Manual:  "set description",
			NamePat: "^(d|desc)$",
			ValPat:  "^.*$",
		},
		{
			Name:    "time",
			Manual:  "set timestamp",
			NamePat: "^(t|time)$",
			ValPat:  "^\\d{4}-\\d{2}-\\d{2};\\d{2}:\\d{2}:\\d{2}$",
		},
	},
}
var altTxnRecordState = &cli.State{
	Action: true,
	Name:   "alt txn record",
	Manual: "alter records of a transaction",
	Pat:    "^record$",
	Param:  2,
}
var delTxnState = &cli.State{
	Action: true,
	Name:   "del txn",
	Manual: "delete some transaction",
	Pat:    "^txn$",
	Param:  1,
}
var newTagState = &cli.State{
	Action: true,
	Name:   "new tag",
	Manual: "create a new tag",
	Pat:    "^tag$",
	Param:  1,
	Flags: []*cli.Flag{
		{
			Name:    "desc",
			Manual:  "set description",
			NamePat: "^(d|desc)$",
			ValPat:  "^.*$",
		},
	},
}
var getTagState = &cli.State{
	Action: true,
	Name:   "get tag",
	Manual: "find a tag",
	Pat:    "^tag$",
}
var getTagByNameState = &cli.State{
	Action: true,
	Name:   "get tag name",
	Manual: "find an tag by name",
	Pat:    "^name$",
	Param:  1,
}
var getTagByDescState = &cli.State{
	Action: true,
	Name:   "get tag desc",
	Manual: "find an tag by desc",
	Pat:    "^desc$",
	Param:  1,
}
var altTagState = &cli.State{
	Action: true,
	Name:   "alt tag",
	Manual: "modify an tag",
	Pat:    "^tag$",
	Param:  1,
	Flags: []*cli.Flag{
		{
			Name:    "desc",
			Manual:  "set description",
			NamePat: "^(d|desc)$",
			ValPat:  "^.*$",
		},
	},
}
var altTagRenameState = &cli.State{
	Action: true,
	Name:   "alt tag rename",
	Manual: "rename a tag",
	Pat:    "^rename$",
	Param:  2,
}
var delTagState = &cli.State{
	Action: true,
	Name:   "del tag",
	Manual: "delete a tag",
	Pat:    "^tag$",
	Param:  1,
}

func init() {
	rootState.SetNext(quitState)
	rootState.SetNext(openState)
	rootState.SetNext(closeState)

	rootState.SetNext(defState)
	rootState.SetNext(calcState)

	rootState.SetNext(newState)
	rootState.SetNext(getState)
	rootState.SetNext(altState)
	rootState.SetNext(delState)

	newState.SetNext(newAccState)

	getState.SetNext(getAccState)
	getAccState.SetNext(getAccChildState)
	getAccState.SetNext(getAccByNameState)
	getAccState.SetNext(getAccByDescState)

	altState.SetNext(altAccState)
	altAccState.SetNext(altAccRenameState)

	delState.SetNext(delAccState)

	newState.SetNext(newTxnState)

	getState.SetNext(getTxnState)
	getTxnState.SetNext(getTxnByIDState)
	getTxnState.SetNext(getTxnByTimeState)
	getTxnState.SetNext(getTxnByDescState)

	altState.SetNext(altTxnState)
	altTxnState.SetNext(altTxnRecordState)

	delState.SetNext(delTxnState)

	newState.SetNext(newTagState)

	getState.SetNext(getTagState)
	getTagState.SetNext(getTagByNameState)
	getTagState.SetNext(getTagByDescState)

	altState.SetNext(altTagState)
	altTagState.SetNext(altTagRenameState)

	delState.SetNext(delTagState)

	parser = cli.NewParser(rootState)
}
