'use client'

import React, { useEffect, useState } from 'react'
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'

import backendUrl from 'lib/backendUrl'

export default function ListDataFiles() {
  interface DataFile {
    CID: string;
    WalletAddress: string;
    Filename: string;
  }

  const [datafiles, setDataFiles] = useState<DataFile[]>([]);

  useEffect(() => {
    fetch(`${backendUrl()}/datafiles`)
      .then(response => {
        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}`);
        }
        return response.json();
      })
      .then(data => {
        console.log('Fetched datafiles:', data);
        setDataFiles(data);
      })
  }, [])

  return (
    <TableContainer>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>CID</TableCell>
            <TableCell>Uploader Wallet</TableCell>
            <TableCell>Filename</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {datafiles.map((datafile, index) => (
            <TableRow key={index}>
              <TableCell>
                <a href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}${datafile.CID}/`}>
                  {datafile.CID}
                </a>
              </TableCell>
              <TableCell>{datafile.WalletAddress}</TableCell>
              <TableCell>{datafile.Filename}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}