package dto

type Paper struct {
	Number  int    // Unique identifier
	Author  string // Author(s) name
	Title   string // Title of the paper
	Format  string // PDF or DOC
	Content []byte // Binary content of the paper
}

type AddPaperArgs struct {
	Author  string
	Title   string
	Format  string
	Content []byte
}

type AddPaperReply struct {
	PaperNumber int
}

type ListPapersReply struct {
	Papers []Paper
}

type GetPaperArgs struct {
	Number int
}

type GetPaperDetailsReply struct {
	Author string
	Title  string
}

type FetchPaperArgs struct {
	Number int
}

type FetchPaperReply struct {
	Content []byte
}
