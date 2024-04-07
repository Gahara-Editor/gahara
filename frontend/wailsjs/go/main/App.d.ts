// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {menu} from '../models';
import {video} from '../models';
import {main} from '../models';

export function AppMenu(arg1:Array<menu.MenuItem>):Promise<menu.Menu>;

export function CreateProjectWorkspace(arg1:string):Promise<string>;

export function DeleteProject(arg1:string):Promise<void>;

export function DeleteProjectFile(arg1:string):Promise<void>;

export function DeleteRIDReferences(arg1:string):Promise<void>;

export function EnableVideoMenus():Promise<void>;

export function ExportVideo(arg1:video.ProcessingOpts):Promise<string>;

export function FilePicker():Promise<void>;

export function GenerateThumbnail(arg1:string):Promise<void>;

export function GetOutputFileSavePath():Promise<string>;

export function GetProjectThumbnail(arg1:string):Promise<string>;

export function GetThumbnail(arg1:string):Promise<void>;

export function GetTimeline():Promise<video.Timeline>;

export function GetTrackDuration():Promise<number>;

export function InsertInterval(arg1:string,arg2:number,arg3:number,arg4:number):Promise<video.VideoNode>;

export function LoadTimeline():Promise<video.Timeline>;

export function ReadGaharaWorkspace():Promise<Array<string>>;

export function ReadProjectWorkspace():Promise<Array<main.Video>>;

export function RemoveInterval(arg1:number):Promise<void>;

export function RenameVideoNode(arg1:number,arg2:string):Promise<void>;

export function ResetTimeline():Promise<void>;

export function SaveTimeline():Promise<void>;

export function SetDefaultAppMenu():Promise<void>;

export function SetProjectDirectory(arg1:string):Promise<void>;

export function SplitInterval(arg1:string,arg2:number,arg3:number,arg4:number):Promise<Array<video.VideoNode>>;
