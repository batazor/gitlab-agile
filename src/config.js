// Set basic configuration app =================================================
export const APP_NAME = process.env.APP_NAME || 'LMap'
export const PORT = process.env.PORT ? parseInt(process.env.PORT) : 4000 // eslint-disable-line no-magic-numbers
export const DOMAIN = process.env.DAMAIN || `http://127.0.0.1:${ PORT }`
export const ENV = process.env.NODE_ENV || 'develop'
export const DEBUG = process.env.NODE_ENV === 'develop'
