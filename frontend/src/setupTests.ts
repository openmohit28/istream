import '@testing-library/jest-dom/vitest'
import { afterEach, beforeEach } from 'vitest'
import { cleanup } from '@testing-library/react'

// Deterministic in-memory localStorage: Node 22+ ships a global localStorage
// that shadows jsdom's and lacks a working implementation under Vitest.
class MemoryStorage implements Storage {
  private store = new Map<string, string>()

  get length() {
    return this.store.size
  }

  clear() {
    this.store.clear()
  }

  getItem(key: string) {
    return this.store.get(key) ?? null
  }

  key(index: number) {
    return [...this.store.keys()][index] ?? null
  }

  removeItem(key: string) {
    this.store.delete(key)
  }

  setItem(key: string, value: string) {
    this.store.set(key, String(value))
  }
}

Object.defineProperty(window, 'localStorage', {
  value: new MemoryStorage(),
  writable: true,
})

beforeEach(() => {
  window.localStorage.clear()
})

afterEach(() => {
  cleanup()
})