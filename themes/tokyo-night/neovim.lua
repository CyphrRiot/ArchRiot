-- Tokyo Night theme for Neovim
-- Authentic Tokyo Night color scheme

return {
  "folke/tokyonight.nvim",
  lazy = false,
  priority = 1000,
  opts = {
    style = "night", -- The theme comes in three styles, `storm`, `moon`, a darker variant `night` and `day`
    transparent = false,
    terminal_colors = false,
    styles = {
      comments = { italic = true },
      keywords = { italic = true },
      functions = {},
      variables = {},
      sidebars = "dark",
      floats = "dark",
    },
    sidebars = { "qf", "help", "neo-tree" },
    day_brightness = 0.3,
    hide_inactive_statusline = false,
    dim_inactive = false,
    lualine_bold = false,
    on_colors = function(colors)
      -- Match Ghostty's exact black background
      colors.bg = "#0a0b10"           -- Primary background (matches Ghostty)
      colors.bg_dark = "#0a0b10"      -- Darker background
      colors.bg_float = "#0a0b10"     -- Float background
      colors.bg_popup = "#0a0b10"     -- Popup background
      colors.bg_sidebar = "#0a0b10"   -- Sidebar background
      colors.bg_statusline = "#0a0b10" -- Statusline background
      colors.fg = "#c0caf5"           -- Foreground text
      colors.fg_dark = "#a9b1d6"      -- Darker foreground
      colors.fg_gutter = "#3b4261"    -- Gutter foreground
      colors.fg_sidebar = "#c0caf5"   -- Sidebar foreground

      -- Tokyo Night accent colors
      colors.blue = "#7aa2f7"         -- Primary blue
      colors.cyan = "#7dcfff"         -- Cyan
      colors.green = "#9ece6a"        -- Green
      colors.magenta = "#bb9af7"      -- Magenta/purple
      colors.orange = "#ff9e64"       -- Orange
      colors.purple = "#9d7cd8"       -- Purple
      colors.red = "#f7768e"          -- Red
      colors.yellow = "#e0af68"       -- Yellow

      -- Additional Tokyo Night colors
      colors.blue0 = "#3d59a1"        -- Dark blue
      colors.blue1 = "#2ac3de"        -- Light blue
      colors.blue2 = "#0db9d7"        -- Bright cyan
      colors.blue5 = "#89ddff"        -- Very light blue
      colors.blue6 = "#b4f9f8"        -- Pale cyan
      colors.blue7 = "#394b70"        -- Dark blue accent
      colors.green1 = "#73daca"       -- Light green
      colors.green2 = "#41a6b5"       -- Teal
      colors.magenta2 = "#ff007c"     -- Hot pink
      colors.purple = "#9d7cd8"       -- Purple
      colors.red1 = "#db4b4b"         -- Dark red
      colors.teal = "#1abc9c"         -- Teal

      -- Git colors
      colors.git = colors.git or {}
      colors.git.add = "#449dab"
      colors.git.change = "#6183bb"
      colors.git.delete = "#914c54"
      colors.gitSigns = colors.gitSigns or {}
      colors.gitSigns.add = "#266d6a"
      colors.gitSigns.change = "#536c9e"
      colors.gitSigns.delete = "#b2555b"
    end,
    on_highlights = function(highlights, colors)
      -- Match Ghostty's exact black background
      highlights.Normal = { bg = "#0a0b10", fg = colors.fg }
      highlights.NormalFloat = { bg = "#0a0b10", fg = colors.fg }
      highlights.NormalNC = { bg = "#0a0b10", fg = colors.fg_dark }
      highlights.CursorLine = { bg = "#292e42" }
      highlights.Visual = { bg = "#33467c" }
      highlights.Search = { bg = colors.orange, fg = colors.bg }
      highlights.IncSearch = { bg = colors.magenta, fg = colors.bg }
      highlights.LineNr = { fg = colors.fg_gutter }
      highlights.CursorLineNr = { fg = colors.orange }

      -- Syntax highlighting
      highlights.Comment = { fg = colors.fg_gutter, italic = true }
      highlights.Keyword = { fg = colors.purple, italic = true }
      highlights.Function = { fg = colors.blue }
      highlights.String = { fg = colors.green }
      highlights.Number = { fg = colors.orange }
      highlights.Boolean = { fg = colors.orange }
      highlights.Constant = { fg = colors.orange }
      highlights.Variable = { fg = colors.fg }
      highlights.Type = { fg = colors.blue1 }
      highlights.Operator = { fg = colors.blue5 }
      highlights.Identifier = { fg = colors.magenta }

      -- UI elements - Match Ghostty black background
      highlights.Pmenu = { bg = "#0a0b10", fg = colors.fg }
      highlights.PmenuSel = { bg = colors.blue7, fg = colors.fg }
      highlights.PmenuSbar = { bg = "#0a0b10" }
      highlights.PmenuThumb = { bg = colors.fg_gutter }
      highlights.SignColumn = { bg = "#0a0b10" }
      highlights.FoldColumn = { bg = "#0a0b10" }

      -- Status line
      highlights.StatusLine = { bg = colors.bg_statusline, fg = colors.fg }
      highlights.StatusLineNC = { bg = colors.bg_statusline, fg = colors.fg_dark }

      -- Tabs
      highlights.TabLine = { bg = colors.bg_statusline, fg = colors.fg_dark }
      highlights.TabLineFill = { bg = colors.bg_statusline }
      highlights.TabLineSel = { bg = colors.bg, fg = colors.fg }
    end,
  },
  config = function(_, opts)
    require("tokyonight").setup(opts)
    vim.cmd([[colorscheme tokyonight-night]])
  end,
}
