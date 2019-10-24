// Copyright 2019 Job Stoit. All rights reserved.

package template

// CreateEntityFrameworkContext creates a dotnet Microsoft.EnittyFrameworkCoreContext


var efContextTmpl = `using Microsoft.EntityFrameworkCore;

namespace {{title .Pkg}}
{
	class Db : DbContext
	{
		public Db(DbContextOptions<Db> options) : base(options)
		{
		}
		{{range .Table}}
		public DbSet<{{title .}}> {{title .}}s { get; set; }{{end}}
	}
	{{range .Table}}
	class {{title .}}
	{
		{{range column .}}{{end}}
	}{{end}}
}`
