$ ->
  searchInput = $('input[type=search]')
  searchForm = $("#search-form")
  resultsTableBody = $('#results table tbody')
  indicesTableBody = $('#indices table tbody')
  urlInput = $('input[type=url]')
  addURLForm = $('#add-url-form')
  addURLButton = $('#add-url')
  listURLsButton = $('#list-url')
  loadMoreButton = $('#load-more')

  # Hide table
  indicesTable = $('#indices')
  resultsTable = $('#results')

  # Vertically center header
  header = $('#header')
  header.css('margin-top', $(window).height() / 2 - header.height() / 1.5)

  # Listen to input event and send search query
  searchInput.on "input", (e) ->
    value = e.target.value
    itemIndex.search value if value

  loadMoreButton.on "click", (e) ->
    itemIndex.loadMore()
    e.preventDefault()

  searchForm.on "submit", (e) ->
    # force search
    searchInput.trigger("input")
    e.preventDefault()

  # When click on radio button trigger search
  $('input[name="search-type"]').on "click", (e) ->
    searchInput.trigger("input")

  addURLForm.on "submit", (e) ->
    addIndexOf urlInput.val() #if urlInput.val()
    e.preventDefault()

  disableClassIn = (elem, clazz, timeout) ->
    setTimeout () ->
      elem.removeClass(clazz)
    , timeout

  insertIntoTableBody = (tbody, indexItem, fields) ->
    row = $("<tr>").addClass("new-entry")
    disableClassIn(row, "new-entry", 2000)
    fields.forEach (field) ->
      td = $("<td>")
      if field == "last_modified_at"
        td.html moment(indexItem[field]).fromNow()
      else if field == "url"
        item = $("<a>").attr("href", indexItem[field]).text(indexItem[field])
        td.append item
      else if field == "name"
        item = $("<a>").attr("href", indexItem["url"]).text(indexItem[field])
        td.append item
      else if field == "size"
        if indexItem[field] == -1 then td.html "- Directory -" else td.html indexItem[field]
      else
        td.html indexItem[field]
      row.append td
    tbody.append row

  listURLsButton.on "click", (e) ->
    realTimeIndexesOf.toggle()

  realTimeIndexesOf = (() ->
    tickHandle: 0

    enable: () ->
      return if @tickHandle != 0
      @update()
      @tickHandle = setInterval () =>
        @update()
      , 1000

    toggle: () ->
      if @tickHandle != 0
        @disable()
      else
        @enable()

    disable: (cb) ->
      if @tickHandle != 0
        clearInterval @tickHandle
      @tickHandle = 0
      indicesTable.fadeOut cb

    update: ->
      $.ajax
        type: "GET"
        url: '/api/indices'
        dataType: 'json'
        success: (data) ->
          indicesTableBody.empty()
          data.forEach (elem) ->
            url = elem.scheme + "://" + elem.host + elem.path
            insertIntoTableBody indicesTableBody, {url: url, count: elem.count}, ["url", "count"]
          header.animate
            'margin-top': 0
            'slow'
          resultsTable.fadeOut ->
            indicesTable.fadeIn()
  )()

  # Send search query with a delay
  itemIndex = (() ->
    from: 0
    timeoutHandle: 0
    query: ""
    step: 10

    loadMore: ->
      @from += @step
      @_loadItems @query, (items) =>
        loadMoreButton.hide() if items.length < @step
        @_addItemsToTable items, () ->
          $("#about-link").ScrollTo
            duration: 2000

    search: (query) ->
      if query != @query
        @query = query
        @from = 0

      # Cancel waiting search query
      if @timeoutHandle != 0
        clearTimeout @timeoutHandle
        @timeoutHandle = 0

      @timeoutHandle = setTimeout () =>
        @_loadItems query, (items) =>
          if items.length >= @step then loadMoreButton.show() else loadMoreButton.hide()
          resultsTableBody.empty()
          @_addItemsToTable items
        @timeoutHandle = 0
      , 300

    _loadItems: (query, cb) ->
      type = $('input[name="search-type"]:checked').val()
      $.ajax
        type: "GET"
        url: '/api/search'
        dataType: 'json'
        data: {q: @query, from: @from, t: type}
        success: cb

    _addItemsToTable: (items, cb) ->
      items.forEach (elem) ->
        insertIntoTableBody resultsTableBody, elem, ['name', 'last_modified_at', 'size']
      header.animate
        'margin-top': 0
        'slow'
      realTimeIndexesOf.disable ->
        resultsTable.fadeIn()
      cb() if cb
  )()

  addIndexOf = (url) ->
    $.ajax
      type: "POST"
      contentType: "application/json; charset=utf-8"
      dataType: 'json'
      url: '/api/indices'
      data: JSON.stringify url: url
      success: (data) ->
        $("#btn-row .error").addClass("hidden")
        $("#btn-row .error").val("")
        urlInput.val("")
        realTimeIndexesOf.enable()
        addURLButton.addClass('green')
        setTimeout () ->
          addURLButton.removeClass('green')
        , 2000
      error: (response, error) ->
        errSpan = $("#btn-row .error")
        addURLButton.addClass('red')
        setTimeout () ->
          errSpan.empty()
          addURLButton.removeClass('red')
        , 2000
        if response.status == 500
          errSpan.html("Internal Error, sorry :-/")
        else if response.status == 422
          errors = response.responseJSON.errors
          out = "<pre>"
          Object.keys(errors).forEach (key) ->
            errors[key].forEach (err) ->
              out += key + ": " + err + "\n"
          errSpan.html(out + "</pre>")
