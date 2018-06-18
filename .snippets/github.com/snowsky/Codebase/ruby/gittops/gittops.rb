require 'octokit'

module GitTops
  class GitTopRepos
    def initialize(token, query, topn=5)
      @token = token
      @query = query
      @topn = topn
    end

    def get_repo
      client = Octokit::Client.new(:access_token => @token)
      #puts client.login

      repos = client.search_repositories(@query, options = {:sort => 'stars', :order => 'desc'})

      if @topn != 0
        repos.items.take(@topn)
      else
        repos.items
      end
     end
   end
end




#'3b8fd948d570f6f6ab56986480d2aa9c153b2b98'
