// 生成彩色占位图（解决 picsum.photos 在国内无法访问的问题）
const gradients = [
  'linear-gradient(135deg,#667eea,#764ba2)',
  'linear-gradient(135deg,#f093fb,#f5576c)',
  'linear-gradient(135deg,#4facfe,#00f2fe)',
  'linear-gradient(135deg,#43e97b,#38f9d7)',
  'linear-gradient(135deg,#fa709a,#fee140)',
  'linear-gradient(135deg,#a18cd1,#fbc2eb)',
  'linear-gradient(135deg,#fccb90,#d57eeb)',
  'linear-gradient(135deg,#e0c3fc,#8ec5fc)',
  'linear-gradient(135deg,#f5576c,#ff6f91)',
  'linear-gradient(135deg,#30cfd0,#330867)',
  'linear-gradient(135deg,#a8edea,#fed6e3)',
  'linear-gradient(135deg,#5ee7df,#b490ca)',
  'linear-gradient(135deg,#d299c2,#fef9d7)',
  'linear-gradient(135deg,#fdfcfb,#e2d1c3)',
  'linear-gradient(135deg,#96deda,#50c9c3)',
  'linear-gradient(135deg,#f5f7fa,#c3cfe2)',
  'linear-gradient(135deg,#667eea,#764ba2)',
  'linear-gradient(135deg,#89f7fe,#66a6ff)',
  'linear-gradient(135deg,#fddb92,#d1fdff)',
  'linear-gradient(135deg,#9890e3,#b1f4cf)',
]

const icons = ['🌅', '🌸', '🌊', '🌺', '🦋', '🌈', '✨', '🌟', '🎨', '🎭', '🎪', '🎡', '🏖️', '🏔️', '🌋', '🏝️', '🎠', '🎯', '🎲', '🧩']

let seedCounter = 0

/**
 * 生成一个可靠的占位图URL（基于SVG data URI）
 * 替代 picsum.photos，在国内网络环境下可正常显示
 * @param {number} width - 图片宽度
 * @param {number} height - 图片高度
 * @param {string} seed - 用于确定颜色和图标的一致性种子
 * @returns {string} data:image/svg+xml 格式的图片URL
 */
export function getPlaceholderImage(width = 400, height = 300, seed = '') {
  seedCounter++
  const hash = seed.split('').reduce((acc, c) => acc + c.charCodeAt(0), 0) || seedCounter
  const gradientIndex = hash % gradients.length
  const iconIndex = hash % icons.length

  const gradient = gradients[gradientIndex]
  const icon = icons[iconIndex]

  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
      <defs>
        <linearGradient id="g" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:${gradient.split(',')[0].replace('linear-gradient(135deg,','')}"/>
          <stop offset="100%" style="stop-color:${gradient.split(',')[1].replace(')','')}"/>
        </linearGradient>
      </defs>
      <rect width="${width}" height="${height}" fill="url(#g)"/>
      <text x="${width/2}" y="${height/2}" text-anchor="middle" dominant-baseline="central"
            font-size="${Math.min(width, height) * 0.35}" filter="drop-shadow(0 2px 8px rgba(0,0,0,0.15))">
        ${icon}
      </text>
    </svg>
  `.trim()

  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
}

/**
 * 生成头像占位图
 */
export function getAvatarPlaceholder(username = '用户', size = 100) {
  return getPlaceholderImage(size, size, username)
}