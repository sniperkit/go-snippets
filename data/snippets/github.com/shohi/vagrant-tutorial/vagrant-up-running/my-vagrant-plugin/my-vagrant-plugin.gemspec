
lib = File.expand_path("../lib", __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require "my-vagrant-plugin/version"

Gem::Specification.new do |spec|
  spec.name          = "my-vagrant-plugin"
  spec.version       = MyVagrantPlugin::VERSION
  spec.authors       = ["shohi"]
  spec.email         = ["oshohi@gmail.com"]

  spec.summary       = "my vagrant plugin on tutorial."

  spec.files         = `git ls-files -z`.split("\x0").reject do |f|
    f.match(%r{^(test|spec|features)/})
  end
  spec.bindir        = "exe"
  spec.executables   = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths = ["lib"]

  spec.add_development_dependency "bundler", "~> 1.16"
  spec.add_development_dependency "rake", "~> 10.0"
end
