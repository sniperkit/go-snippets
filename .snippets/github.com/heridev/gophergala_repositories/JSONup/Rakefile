require 'sprockets/standalone'

Sprockets::Standalone::RakeTask.new(:assets) do |task, sprockets|
  task.assets   = %w(css/app.css js/app.js)
  task.sources  = %w(assets)
  task.output   = File.expand_path('../public', __FILE__)
  task.compress = false
  task.digest   = false

  #sprockets.js_compressor  = :uglifier
  #sprockets.css_compressor = :sass
end

task :build => 'assets:compile' do
  html_file = File.expand_path('../assets/index.html', __FILE__)
  public_dir = File.expand_path('../public/', __FILE__)
  `cp #{html_file} #{public_dir}`
end

task :default => :build
