
require 'yaml'
require 'twitter-text'
TwitterExtractor = Twitter::Extractor
require 'twitter'
require 'travis'

module Tweetcode
  class CI
    def initialize
      scheduled = ENV['TWEET_ON_CRON_OR_API_ONLY'] == 'true'
      @mode = case
              when event_type?('push', 'pull_request') && branch != 'master'
                :test
              when ci? && branch == 'master'
                if event_type?('cron', 'api') && scheduled
                  :tweet_scheduled
                elsif event_type?('push') && !scheduled
                  :tweet
                else
                  puts "SKIPPING BUILD"
                  puts "TWEET_ON_CRON_OR_API_ONLY=#{scheduled}"
                  puts "EVENT TYPE IS #{ENV['TRAVIS_EVENT_TYPE']}"
                  exit(0)
                end
              else
                :preview
              end
      @errors = []

      if @mode == :tweet_scheduled
        acc_tok = ENV['TRAVIS_API_ACCESS_TOKEN'] || fatal("TRAVIS_API_ACCESS_TOKEN environment variable is not set")
        client = Travis::Client.new(access_token: acc_tok)
        @travis_env = client.repo(ENV['TRAVIS_REPO_SLUG']).env_vars
      end
    end

    def start
      begin
        send(@mode)
        exit(0)
      rescue => e
        print_errors
        puts e.to_s
        puts e.backtrace
        fatal "FAIL"
      end
    end

    private

    def error(e)
      @errors << e
      nil
    end

    def fatal(e)
      error(e)
      print_errors
      exit(1)
    end

    def print_errors
      puts "ERRORS:"
      puts @errors
    end

    def commits_cnt
      @commits_cnt ||= run("git rev-list --count '#{commit_id}' '^master'").to_i
    end

    def next_commit(id)
      c = run("git rev-list #{id}..HEAD --reverse").lines[0]
      unless c
        puts "Nothing to tweet, exiting..."
        exit(0)
      end
      c.strip
    end

    def skip_commit
      @cid = next_commit(@cid)
      @raw_msg = nil
      @file = nil
    end

    def last_scheduled_commit_id
      @last_scheduled_id ||= ENV['TWEETCODE_LAST_SCHEDULED_ID']
    end

    def commit_id
      @cid ||= lambda do
        if @mode == :tweet_scheduled && last_scheduled_commit_id
          next_commit(last_scheduled_commit_id)
        else
          run("git rev-list --no-merges -n 1 HEAD").strip
        end
      end.call
    end

    def short_commit_id
      run("git rev-parse --short=4 #{commit_id}").strip
    end

    def commit_msg
      @raw_msg ||= run("git log --format=%B -n 1 #{commit_id}")
    end

    def file
      @file ||= lambda do
        gshow = run("git show --name-only #{commit_id}").lines.reverse
        files = gshow[0..gshow.index{|l| l.match(/^\s*/)}].map(&:strip)
        verbose("Files in commit:\n#{files}\n")

        if files.length != 1
          fatal 'Commit must add only 1 new file'
        elsif !files[0].start_with?('snippets/') || !files[0].end_with?(ENV['SNIPPET_FILE_EXTENSION'] || "")
          error "Snippet file must be in 'snippets' directory and have *.#{ENV['SNIPPET_FILE_EXTENSION']||'*'} extension"
        end

        if File.read(files[0]).lines.length > (ENV['SNIPPET_MAX_LINES'].to_i || 20)
          error "Snippet can't have more than #{ENV['SNIPPET_MAX_LINES'] || 20} lines"
        end
        files[0]
      end.call
    end

    def parsed_msg
      while ['[skip ci]', '[ci skip]'].any?{|t| commit_msg.include?(t)} do
        skip_commit
      end
      begin
        @parsed_msg ||= YAML::load(commit_msg)
        if @parsed_msg['note'].nil? || @parsed_msg['twitter_username'].nil?
          fatal '"note" and "twitter_username" are required fields in commit message'
        end
        @parsed_msg
      rescue
        fatal("Failed to parse commit message (it must be a valid YAML, see .gitmessage reference and README)")
      end
    end

    def twitter_username
      @twitter_username ||= parsed_msg['twitter_username'].tr('@', '')
    end

    def note
      parsed_msg['note']
    end

    def source
      parsed_msg['source']
    end

    def tweet_text
      @tweet_text ||= begin
        m = [note, source ? "Source: #{source}\n" : nil,  "[#{short_commit_id}] by @#{twitter_username}"].compact.join("\n")
        urls = TwitterExtractor.extract_urls(m)
        if tweet_len = (urls.length > 0 ? m.length - urls.map(&:length).reduce{|a,b| a+b} + urls.length * 23 : m.length) > 140
          urls.each { |u| m.gsub(u, "#{u[0..20]}...") }
          error "Generated tweet has #{tweet_len} characters, 140 is maximum allowed (any URL is counted as 23 characters): \n#{m}"
        else
          verbose "Tweet length: #{tweet_len}"
        end
        verbose("Generated tweet text:\n#{m}")
        m
      rescue => e
        fatal("FATAL ERROR: #{e.inspect}\n#{e.backtrace.join("\n")}")
      end
    end

    def validate_snippet
      if cmd = (ENV['VALIDATE_SNIPPET_COMMAND'].tr('"', '').tr("'", '') rescue false)
        ENV['SNIPPET_FILE'] = file
        if !system(cmd)
          error "Snippet validation command failed:\n"
        end
      end
    end

    def make_png
      font = ENV['SNIPPET_FONT_NAME'].tr("'", '').tr('"', '') rescue 'Source Code Pro'
      font_size = ENV['SNIPPET_FONT_SIZE'] || 19
      color_scheme = ENV['SNIPPET_COLOR_SCHEME'].tr("'", '').tr('"', '') rescue 'Tomorrow-Night'
      run("#{File.join(run("npm bin").strip, 'copper')} -f='#{font}' --fontSize=#{font_size} -t=#{color_scheme} #{file}").split.last
    end

    def tweet
      t_txt = tweet_text
      png_path = make_png
      unless File.file?(png_path)
        puts("Failed to generate image:")
        puts png_path
        exit(1)
      end
      t = Twitter::REST::Client.new do |config|
        config.consumer_key        = ENV["TWITTER_CONSUMER_KEY"]
        config.consumer_secret     = ENV["TWITTER_CONSUMER_SECRET"]
        config.access_token        = ENV["TWITTER_ACCESS_TOKEN"]
        config.access_token_secret = ENV["TWITTER_ACCESS_SECRET"]
      end
      begin
        t.update_with_media(t_txt, File.open(png_path, 'r'))
        puts("Tweet has been successfully published")
      rescue => e
        puts "Failed to publish tweet:"
        puts e.to_s
        exit(1)
      ensure
        run("rm -f #{png_path}")
      end
    end

    def tweet_scheduled
      tweet
      verbose("Setting TWEETCODE_LAST_SCHEDULED_ID from #{ENV['TWEETCODE_LAST_SCHEDULED_ID']} to #{commit_id}")
      @travis_env.upsert('TWEETCODE_LAST_SCHEDULED_ID', commit_id, public: true)
    end

    def test
      verbose("Commits in branch: #{commits_cnt}")
      if commits_cnt != 1
        puts 'Branch must have single commit. See contribution flow'
        exit(1)
      end
      tweet_text
      validate_snippet
      if @errors.any?
        fatal("BUILD FAILED")
      end
    end

    def preview
      verbose("Commits in branch: #{commits_cnt}")
      if commits_cnt != 1 && branch != 'master'
        puts 'Branch must have single commit. See contribution flow'
        exit(1)
      end
      set_env_from_config
      png_path = make_png
      run "mv #{png_path} preview.png"
      puts "Tweet text:\n\n#{tweet_text}"
      validate_snippet
      if @errors.any?
        fatal('Preview failed')
      end
      puts "\nImage: preview.png"
    end

    def set_env_from_config
      cfg = YAML::load(File.read('.travis.yml'))
      cfg['env']['global'].each do |e|
        var, val = e.split('=')
        ENV[var] = val
      end
    end

    def pr?
      ENV['TRAVIS_PULL_REQUEST'] != 'false'
    end

    def branch
      @current_branch ||= if ci?
        (pr? ? ENV['TRAVIS_PULL_REQUEST_BRANCH'] : ENV['TRAVIS_BRANCH']) || 'HEAD'
      else
        run("git rev-parse --abbrev-ref HEAD")
      end
    end

    def ci?
      ENV['CI']
    end

    def event_type?(*type)
      type.include? ENV['TRAVIS_EVENT_TYPE']
    end

    VERBOSE = ENV['VERBOSE'] && !['0', 'false', 'FALSE'].include?(ENV['VERBOSE'])

    def verbose(msg)
      if VERBOSE
        puts "#{msg}"
      end
    end

    def run(cmd)
      verbose "Running `#{cmd}`"
      output = `#{cmd.strip}`.strip
      verbose "Output:\n#{output}\n"
      output
    end
  end
end

def ci
  tc = Tweetcode::CI.new
  tc.start
end

ci()
