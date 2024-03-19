export namespace main {
	
	export class Interval {
	    start: number;
	    end: number;
	
	    static createFrom(source: any = {}) {
	        return new Interval(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.start = source["start"];
	        this.end = source["end"];
	    }
	}
	export class Video {
	    name: string;
	    extension: string;
	    filepath: string;
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.extension = source["extension"];
	        this.filepath = source["filepath"];
	    }
	}

}

export namespace video {
	
	export class VideoNode {
	    rid: string;
	    id: string;
	    start: number;
	    end: number;
	
	    static createFrom(source: any = {}) {
	        return new VideoNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rid = source["rid"];
	        this.id = source["id"];
	        this.start = source["start"];
	        this.end = source["end"];
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

