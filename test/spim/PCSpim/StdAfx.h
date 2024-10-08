// SPIM S20 MIPS simulator.
// Definitions for the SPIM S20.
//
// Copyright (c) 1990-2010, James R. Larus.
// Changes for DOS and Windows versions by David A. Carley (dac@cs.wisc.edu)
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// Neither the name of the James R. Larus nor the names of its contributors may
// be used to endorse or promote products derived from this software without
// specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
//

// stdafx.h : include file for standard system include files,
//  or project specific include files that are used frequently, but
//      are changed infrequently
//

#define VC_EXTRALEAN  // Exclude rarely-used stuff from Windows headers

#include <afxwin.h>  // MFC core and standard components
#include <afxext.h>  // MFC extensions
#ifndef _AFX_NO_AFXCMN_SUPPORT
#include <afxcmn.h>  // MFC support for Windows 95 Common Controls
#endif               // _AFX_NO_AFXCMN_SUPPORT

#include <afxtempl.h>  // MFC template classes

#include "..\CPU\spim.h"
#include "..\CPU\string-stream.h"
#include "..\CPU\spim-utils.h"
#include "..\CPU\inst.h"
#include "..\CPU\sym-tbl.h"
#include "..\CPU\reg.h"
#include "..\CPU\mem.h"
#include "..\CPU\scanner.h"

#ifdef STDAFX_CPP
#define GLOBAL
#else
#define GLOBAL extern
#endif

class CPCSpimView;

GLOBAL CPCSpimView *g_pView;
GLOBAL BOOL g_fSaveWinPos;
GLOBAL BOOL g_fLoadExceptionHandler;
GLOBAL BOOL g_fRunning;
GLOBAL BOOL g_fGenRegHex;
GLOBAL BOOL g_fFPRegHex;
GLOBAL BOOL g_checkUndefinedSymbols;
GLOBAL CString g_strCmdLine;
