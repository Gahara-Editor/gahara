export namespace main {
	
	export class Video {
	    id: string;
	    name: string;
	    extension: string;
	    filepath: string;
	    duration: number;
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.extension = source["extension"];
	        this.filepath = source["filepath"];
	        this.duration = source["duration"];
	    }
	}

}

export namespace video {
	
	export class ProcessingOpts {
	    resolution: string;
	    codec: string;
	    crf: string;
	    preset: string;
	    input_path?: string;
	    output_path: string;
	    filename: string;
	    video_format: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessingOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.resolution = source["resolution"];
	        this.codec = source["codec"];
	        this.crf = source["crf"];
	        this.preset = source["preset"];
	        this.input_path = source["input_path"];
	        this.output_path = source["output_path"];
	        this.filename = source["filename"];
	        this.video_format = source["video_format"];
	    }
	}
	export class VideoNode {
	    start: number;
	    end: number;
	    rid: string;
	    id: string;
	    name: string;
	    losslessexport: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VideoNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.start = source["start"];
	        this.end = source["end"];
	        this.rid = source["rid"];
	        this.id = source["id"];
	        this.name = source["name"];
	        this.losslessexport = source["losslessexport"];
	    }
	}
	export class Timeline {
	    video_nodes: VideoNode[];
	
	    static createFrom(source: any = {}) {
	        return new Timeline(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.video_nodes = this.convertValues(source["video_nodes"], VideoNode);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

