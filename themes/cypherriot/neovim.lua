-- CypherRiot theme for Neovim
-- Based on your current color scheme

return {
  "folke/tokyonight.nvim",
  lazy = false,
  priority = 1000,
  opts = {
    style = "night", -- The theme comes in three styles, `storm`, `moon`, a darker variant `night` and `day`
    transparent = false,
    terminal_colors = true,
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
      -- CypherRiot custom colors based on your waybar theme
      colors.bg = "#1a1b26"           -- Primary background
      colors.bg_dark = "#16161e"      -- Darker background
      colors.bg_float = "#222436"     -- Float background (waybar bgcolor)
      colors.bg_popup = "#222436"     -- Popup background
      colors.bg_sidebar = "#1a1b26"   -- Sidebar background
      colors.bg_statusline = "#222436" -- Statusline background
      colors.fg = "#ffffff"           -- Foreground (waybar fgcolor)
      colors.fg_dark = "#c7c7c7"      -- Darker foreground
      colors.fg_gutter = "#686868"    -- Gutter foreground
      colors.fg_sidebar = "#ffffff"   -- Sidebar foreground
      
      -- Accent colors from your waybar
      colors.blue = "#7da6ff"         -- accent5
      colors.cyan = "#0db9d7"         -- accent2
      colors.green = "#9ece6a"        -- From your mic button
      colors.magenta = "#bb9af7"      -- accent4
      colors.orange = "#ff9e64"       -- accent3
      colors.purple = "#9d7bd8"       -- workspace_active
      colors.red = "#ff7a93"          -- accent1
      colors.yellow = "#e0af68"       -- CPU warning color
      
      -- Special colors
      colors.git.add = "#9ece6a"
      colors.git.change = "#ff9e64"
      colors.git.delete = "#ff7a93"
      colors.gitSigns.add = "#9ece6a"
      colors.gitSigns.change = "#ff9e64"
      colors.gitSigns.delete = "#ff7a93"
    end,
    on_highlights = function(highlights, colors)
      -- Custom highlights for CypherRiot theme
      highlights.Normal = { bg = colors.bg, fg = colors.fg }
      highlights.NormalFloat = { bg = colors.bg_float, fg = colors.fg }
      highlights.CursorLine = { bg = "#4a2b7a" } -- window_color from waybar
      highlights.Visual = { bg = "#4a2b7a" }
      highlights.Search = { bg = colors.orange, fg = colors.bg }
      highlights.IncSearch = { bg = colors.magenta, fg = colors.bg }
    end,
  },
  config = function(_, opts)
    require("tokyonight").setup(opts)
    vim.cmd([[colorscheme tokyonight]])
  end,
}
