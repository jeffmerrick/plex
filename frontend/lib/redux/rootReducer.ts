/* Instruments */
import {
  userSlice,
  dataFileAddSlice,
  dataFileListSlice,
  toolAddSlice,
  toolListSlice,
  flowAddSlice,
  flowListSlice,
  flowDetailSlice,
  jobDetailSlice,
} from './slices'


export const reducer = {
  user: userSlice.reducer,
  dataFileAdd: dataFileAddSlice.reducer,
  dataFileList: dataFileListSlice.reducer,
  toolAdd: toolAddSlice.reducer,
  toolList: toolListSlice.reducer,
  flowAdd: flowAddSlice.reducer,
  flowList: flowListSlice.reducer,
  flowDetail: flowDetailSlice.reducer,
  jobDetail: jobDetailSlice.reducer,
}
