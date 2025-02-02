'use client'

import React from 'react'
import { useRouter } from 'next/navigation'

export function useAuth() {
  const user = localStorage.getItem('username')
  const walletAddress = localStorage.getItem('walletAddress')
  const router = useRouter();

  React.useEffect(() => {
    if (!user || !walletAddress) {
      router.push('/login')
    }
  }, [user])
}
