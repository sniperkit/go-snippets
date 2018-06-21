1. cargo new --bin hello_world: 
	a. cargo new 在当前目录下新建一个基于cargo项目管理的rust项目，项目名称为hello_world，--bin表示该项目将生成可执行文件。
2. cargo build 
3. cargo run
4. cargo项目的典型结构：
	a. 源码位于src目录下。
	b. 默认库的入口位于src/lib.rs。
	c. 默认的可执行程序入口文件是src/main.rs.
	d. 其他可选的可执行文件位于src/bin/*.rs（这里每一个rs文件对应一个可执行程序）。
	e. 外部测试源代码文件位于tests目录下。
	f. 示例程序源码位于examples目录下。
	g. 基准测试源码位于benches目录下.
5. Cargo.toml 和Cargo.lock是cargo项目代码管理的核心文件，cargo工具的所有活动均基于这两个文件。
6. Cargo.toml是cargo特有的项目数据描述文件,cargo.toml文件存储了项目的所有信息，它直接面向程序员。如果想让rust项目以正确的方式进行构建，测试和运行，那么必须按照合理的方式构建cargo.toml.
7. Cargo.lock文件则不直接面向程序员, 一般情况下不需要直接修改该文件。Cargo.lock文件是cargo工具根据统一目录的Cargo.toml文件生成的项目依赖详细清单。
8. Cargo.toml文件是由诸如[package] [dependencies]这样的段落组成，每一个短路又由多个字段组成，这些段落和字段就描述了该项目组织的基本信息。
	1. package段落： package段落描述了软件开发者对本项目的各种元数据描述信息,例如项目名称(name)，版本(version)，作者(authers)等。
	2. dependencies段落：cargo的最大优势在于，能够对该项目的各种依赖进行方便，同一和灵活的管理。在Cargo.toml中主要通过各种依赖段落来描述该项目的所有依赖信息。Cargo.toml中常用的依赖段落包括：
		a. 基于rust官方仓库crates.io，通过版本说明来描述。
		b. 基于项目源代码的git仓库地址，通过URL来描述。
		c. 基于本地项目的绝对路径或者相对路径，通过类Unix模式的路径来描述。

		eg:
			[dependencies]
			typemap="0.3"
			plugin="0.2"
			hamer={version = "0.5.0"}
			color={git = "https://github.com/bjz/color-rs"}
			geometry={path = "crates/geometry"}
	3. test段落：rust中单元测试主要通过在项目代码的测试代码部分使用#[test]属性来描述，而集成测试则一般都会通过Cargo.toml文件中的[test]段落来进行描述。
		eg: 
			[[test]]
			name = "testinit"
			path = "tests/testinit.rs"

			[[test]]
			name = "testtime"
			path = "tests/testtime.rs"
		其中name字段定义的集成测试的名称，path字段定义了集成测试文件相对于toml文件的路径。
	4. example段落：对于example段落中声明的examples，需要通过cargo run --examples NAME来运行，其中NAME对应于你在name字段中定义的名称.
		eg:
			[[example]]
			name = "timeout"
			path = "example/timeout.rs"
	5. bin段落： 对于bin段落中声明的bin，需要通过cargo run --bin NAME来运行，其中NAME对应于你在name字段中定义的名称.
		eg:
			[[bin]]
			name = "bin1"
			path = "bin/bin1/rs"
9. cargo clean
10. cargo install
